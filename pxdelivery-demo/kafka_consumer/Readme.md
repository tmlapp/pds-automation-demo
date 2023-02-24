# Kafka Consumer

This directory contains golang code for reading orders from a Kafka queue and publishing them to a mysql database.

> If the schema for the order submissions change in px-delivery, the Kafka consumer will need to change as well. 