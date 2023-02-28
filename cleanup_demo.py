import pds_rest
import json
import os
from pathlib import Path

pathMySQL = 'my-my-ds.json'
pathMongo = 'mdb-my-ds.json'
pathKafka = 'kf-my-ds.json'

mySqlObj = Path(pathMySQL)
mongoObj = Path(pathMongo)
kafkaObj = Path(pathKafka)

os.system('helm -n px-delivery uninstall pxdeliver')

if mySqlObj.exists():
     with open('my-my-ds.json', 'r') as openfile:
         x = json.load(openfile)
         ds_id = json.dumps(x["id"])
     print("MySQL ID is:" + ds_id)
     pds_rest.delete_deployments(ds_id)

if mongoObj.exists():
      with open('mdb-my-ds.json', 'r') as openfile:
          x = json.load(openfile)
          ds_id = json.dumps(x["id"])
      print("Mongo ID is:" + ds_id)
      pds_rest.delete_deployments(ds_id)


if kafkaObj.exists():
      with open('kf-my-ds.json', 'r') as openfile:
          x = json.load(openfile)
          ds_id = json.dumps(x["id"])
      print("Kafka ID is:" + ds_id)
      pds_rest.delete_deployments(ds_id)