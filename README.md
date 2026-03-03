# Kart Statistics (Currently In Development)
## Description
- This is essentially a web scraping project for data analytics for the local go karting track (Redline Racing - Orem).
- This will have a few different components:
  - A cron job that grabs kart data using change data capture.
  - A database that holds kart data.
  - A detailed stat board for my friends and I to track times and averages.
    - Total average per person
    - Average per person for last 10 heats
    - Total track record per person
    - Track record for last 10 heats
    - Kart specific statistics
      - total average time per kart
      - average time per kart for last 10 heats
      - total track records
      - track record for last 10 heats
  - Some sort of integration that will post messages to Discord displaying the leaderboard.

___

## Resources
- [DB Schema](https://drawsql.app/teams/student-1578/diagrams/kart-stats)

___

## Steps
- [x] Make DB Schema for scrapeable data
- [x] Make DB initialization script
- [ ] Initialize local DB
- [ ] Create web scraper to populate DB (using Change Data Capture)
- [ ] Test web scraper with local DB
- [ ] Research cloud DB options
- [ ] Implement cloud DB
- [ ] Set up a cron job that uses the web scraper to populate cloud DB
- [ ] Grant access to cloud database to other contributors
- [ ] Make a script that determines a leader-board
- [ ] Research discord integrations for displaying leader-board changes
- [ ] Implement discord integration for displaying leader-board changes
- [ ] Plan out a possible website representation

___

## Important to Note
- The data listed is all public and can be queried without authorization, so there is not a privacy issue here.

