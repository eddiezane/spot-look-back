# Spot Look Back

This is a simple app that pulls your 20 most recently played tracks from Spotify and stores them in a PostgreSQL database for future querying. I run this from a Raspberry Pi on a cronjob every 10 minutes.

## Usage

```bash
*/10 * * * * /opt/spot-look-back/spot-look-back -db "postgresql://username:password@host:port/spot-look-back?sslmode=require" -token "Spotify Refresh Token" -clientID "Spotify App Client ID" -clientSecret "Spotify App Client Secret" >> /opt/spot-look-back/out.log
```
