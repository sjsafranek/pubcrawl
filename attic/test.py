

import requests

import conf


url = 'https://api.foursquare.com/v2/venues/search'
resp = requests.get(url, params = {
	'client_id': conf.CLIENT_ID,
	'client_secret': conf.CLIENT_SECRET,
	'v': '20191212',
	'll': '40,-74'
})

print(resp.text)
