
# pip3 install foursquare

# /api/v1/crawl?longitude=-123.088000&latitude=44.046174
import uuid
import random
import foursquare
from flask import Flask
from flask import request
from flask import jsonify
from numpy.random import choice

import conf
from cache import DiskCache


cache = DiskCache(directory="data")

app = Flask(__name__)


# Construct the client object
client = foursquare.Foursquare(
    client_id = conf.CLIENT_ID,
    client_secret = conf.CLIENT_SECRET
)

def newPlacesCrawl(longitude, latitude, name=None):
    results = client.venues.search(
        params = {
            # 'query': 'bar',
            'limit': 50,
            'radius': 2000,     # meters
            'categoryId': ','.join(list(conf.CATEGORIES.values())),
            'll': '{0},{1}'.format(latitude, longitude)
        })
    places = results['venues']
    votes = {}
    visited = {}
    for place in places:
        votes[place['id']] = 1
        visited[place['id']] = False
    return {
        "id": str(uuid.uuid4()),
        "name": name,
        "votes": votes,
        "visited": visited,
        "places": places
    }


@app.route('/api/v1/crawl')
def buildPlacesCrawl():
    longitude = request.args.get("longitude", None)
    latitude = request.args.get("latitude", None)
    name = request.args.get("name", None)
    if not longitude or not latitude:
        return jsonify({"status":"error","message":"bad_request"})
    crawl = newPlacesCrawl(longitude, latitude, name=name)
    cache.set(crawl['id'], crawl)
    return jsonify({
        "status": "ok",
        "data": {
            "crawl": crawl
        }
    })


@app.route('/api/v1/crawl/<crawl_id>')
def fetchPlacesCrawl(crawl_id):
    # filter on get...
    n = request.args.get("n", None)
    if n:
        # random.shuffle(pubs)
        print("TODO!!!")
    crawl = cache.get(crawl_id)
    if crawl:
        return jsonify({
            "status": "ok",
            "data": {
                "crawl": crawl
            }
        })
    return jsonify({"status":"error","message":"not_found"})


@app.route('/api/v1/crawl/<crawl_id>/vote/<place_id>')
def votePlaceCrawlPlace(crawl_id, place_id):
    crawl = cache.get(crawl_id)
    if crawl:
        for place in crawl['places']:
            if place['id'] == place_id:
                crawl['votes'][place_id] += 1
                cache.set(crawl['id'], crawl)
                return jsonify({"status": "ok"})
    return jsonify({"status":"error","message":"not_zfound"})


@app.route('/api/v1/crawl/<crawl_id>/visited/<place_id>')
def visitPlaceCrawlPlace(crawl_id, place_id):
    crawl = cache.get(crawl_id)
    if crawl:
        for place in crawl['places']:
            if place['id'] == place_id:
                crawl['visited'][place_id] = True
                cache.set(crawl['id'], crawl)
                return jsonify({"status": "ok"})
    return jsonify({"status":"error","message":"not_found"})


@app.route('/api/v1/crawl/<crawl_id>/next')
def nextPlaceCrawlPlace(crawl_id):
    crawl = cache.get(crawl_id)
    if crawl:

        number_of_items_to_pick = 1
        list_of_candidates = []
        probability_distribution = []

        total_votes = 0
        for place in crawl['visited']:
            if not crawl['visited'][place]:
                list_of_candidates.append(place)
                total_votes += crawl['votes'][place]

        for place in list_of_candidates:
            probability_distribution.append(
                crawl['votes'][place] / total_votes
            )

        print('list_of_candidates', list_of_candidates)
        print('number_of_items_to_pick', number_of_items_to_pick)
        print('probability_distribution', probability_distribution)

        placeIds = choice(
                    list_of_candidates,
                    number_of_items_to_pick,
                    p = probability_distribution)
        placeId = placeIds[0]

        selectedPlace = {}
        for place in crawl['places']:
            if place['id'] == placeId:
                selectedPlace = place
                break

        return jsonify({"status":"ok","place": selectedPlace})
    return jsonify({"status":"error","message":"not_found"})

#
