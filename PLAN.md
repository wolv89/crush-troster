There is a tennis club that participates in local competition throughout the year. The competition is manage through a system called TROLS. At the start of a given season the fixture is published for all of the teams. I would like to you to build a simple web app for downloading the fixtures for all of the teams within a given club and competition, and allocating them to the available courts within the club.

Specifically:
This system will be built with Golang, using just it's standard library. It will launch a localhost web app which will require some more complex interactivity. This can be handled with any of vanilla Javascript, jQuery, React, etc - but it will only be loaded via the Go app, not over a separate runtime like Node or Bun.

TROLS:
https://trols.org.au/brta/fixture.php
This is not an API, this is just a website which is fairly old. The layout is built using HTML tables, and the navigation for this page is done by selecting an item from the dropdown/select element, which triggers a html form submission!

The breakdown will look like this:

1. Query the user for the name of their tennis club, how many courts they have, and how many courts are needed for each home team, in the terminal, when launching
2. Pull a list of available competitions from the select field in the trols URL above
3. Confirm which competition the user wants to manage
4. Step through each subnavigation option within that competition on Trols, and download any pages that belong to the club mentioned in step 1.
	a. As mentioned, this is not an API endpoint, so add a small delay (~500ms) between hitting these pages, to not overwhelm the host
	b. Suggest saving the raw HTML to memory, or even to a temp dir, while running
5. Parse the captured HTML pages to pull out the details of which teams are playing each week, and whether they are the home or away team.
	a. Clubs can have multiple teams in the same competition, who may be playing each other, with one of them designated the home team
6. Launch a web app using the HTTP library
7. Build a table like listing of the dates and which teams are playing at home each week, across all the competitions searched above.

LATER:
We will need to assign teams to different courts in a roster, with options to shuffle them, and track how many times they play on each court (whether as the home or away team) so we can craft a balanced roster.

