import os.path
import json

# @object DiskCache
class DiskCache(object):

    def __init__(self, directory=""):
        self.directory = directory

    def get(self, key):
        filepath = '{0}/{1}.json'.format(self.directory, key)
        if os.path.exists(filepath):
            with open(filepath, 'r') as fileHandler:
                return json.load(fileHandler)
        return None

    def set(self, key, value={}):
        filepath = '{0}/{1}.json'.format(self.directory, key)
        with open(filepath, 'w') as fileHandler:
            json.dump(value, fileHandler)
