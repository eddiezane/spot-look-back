package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/zmb3/spotify"
	"golang.org/x/oauth2"

	log "github.com/sirupsen/logrus"

	"github.com/jackc/pgx/pgtype"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
)

var (
	dbConnectionString = flag.String("db", "", "The PostgreSQL database connection string: postgresql://username:password@host:port/database?sslmode=require")
	refreshToken       = flag.String("token", "", "Your Spotify refresh token. Can be obtained https://developer.spotify.com/documentation/web-api/quick-start/")
	clientID           = flag.String("clientID", "", "Your Spotify Application Client ID")
	clientSecret       = flag.String("clientSecret", "", "Your Spotify Application Client Secret")
)

type RecentlyPlayedTrack struct {
	UserID      string            `db:"user_id"`
	PlayedAt    time.Time         `db:"played_at"`
	Duration    int               `db:"duration"`
	TrackID     string            `db:"track_id"`
	TrackName   string            `db:"track_name"`
	ArtistIDs   *pgtype.TextArray `db:"artist_ids"`
	ArtistNames *pgtype.TextArray `db:"artist_names"`
}

func (r *RecentlyPlayedTrack) String() string {
	var artistIDs []string
	var artistNames []string
	r.ArtistIDs.AssignTo(&artistIDs)
	r.ArtistNames.AssignTo(&artistNames)

	return fmt.Sprintf("user_id: %s played_at: %v duration: %d track_id: %s track_name: %s artist_ids: %s artist_names: %s", r.UserID, r.PlayedAt, r.Duration, r.TrackID, r.TrackName, strings.Join(artistIDs, ", "), strings.Join(artistNames, ", "))
}

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	flag.Parse()
}

func main() {
	if *dbConnectionString == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}
	if *refreshToken == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}
	if *clientID == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}
	if *clientSecret == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	log.Info("spot-look-back starting...")

	db, err := sqlx.Connect("pgx", *dbConnectionString)
	defer db.Close()
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	token := &oauth2.Token{
		RefreshToken: *refreshToken,
	}
	auth := spotify.NewAuthenticator("")
	auth.SetAuthInfo(*clientID, *clientSecret)
	client := auth.NewClient(token)

	user, err := client.CurrentUser()
	if err != nil {
		log.Fatal(err)
	}
	userID := user.ID

	tracks, err := client.PlayerRecentlyPlayed()
	if err != nil {
		log.Fatal(err)
	}

	rows := make([]RecentlyPlayedTrack, 0, len(tracks))
	for _, track := range tracks {
		artists := track.Track.Artists
		artistIDs := make([]string, 0, len(artists))
		artistNames := make([]string, 0, len(artists))
		for _, artist := range artists {
			artistIDs = append(artistIDs, artist.ID.String())
			artistNames = append(artistNames, artist.Name)
		}
		idArray := &pgtype.TextArray{}
		idArray.Set(artistIDs)
		nameArray := &pgtype.TextArray{}
		nameArray.Set(artistNames)
		r := &RecentlyPlayedTrack{
			UserID:      userID,
			PlayedAt:    track.PlayedAt,
			Duration:    track.Track.Duration,
			TrackID:     track.Track.ID.String(),
			TrackName:   track.Track.Name,
			ArtistIDs:   idArray,
			ArtistNames: nameArray,
		}
		rows = append(rows, *r)
	}

	insertStmt := "INSERT INTO tracks VALUES(:user_id, :played_at, :duration, :track_id, :track_name, :artist_ids, :artist_names) ON CONFLICT DO NOTHING"
	for _, r := range rows {
		result, err := db.NamedExec(insertStmt, r)
		if err != nil {
			log.Fatal(err)
		}
		rowsAffected, err := result.RowsAffected()
		if err != nil {
			log.Fatal(err)
		}
		if rowsAffected > 0 {
			var artistIDs []string
			var artistNames []string
			r.ArtistIDs.AssignTo(&artistIDs)
			r.ArtistNames.AssignTo(&artistNames)

			log.WithFields(log.Fields{
				"user_id":      r.UserID,
				"played_at":    r.PlayedAt,
				"duration":     r.Duration,
				"track_id":     r.TrackID,
				"track_name":   r.TrackName,
				"artist_ids":   artistIDs,
				"artist_names": artistNames,
			}).Info("row added")
		}
	}

	log.Info("spot-look-back done...")
}
