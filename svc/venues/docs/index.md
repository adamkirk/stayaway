# Venues

Contains all the venues for the wider system.

## Models

### organisation (this may wanna be separate service, but in here is fine for now)

- id
- name
- slug

### Venue

- id
- chain (nullable)
- name
- slug
- type
- address_line_1
- address_line_2
- conurbation
- county
- post_code
- available_facilities

### Room Spec

- id
- venue_id
- name
- type (shared|private)
- min_occupancy
- max_occupancy
- checkin
- early_checkin
- checkout
- late_checkout
- standard_price
- description
- room_number_generation_rules
- room_allocation_strategy
- room_facilities
- availability

### Conurbation

See: https://simplemaps.com/data/gb-cities

> It ain't a full list, but it'll be enough...dont' wanna pay for the full list tbh

- id
- name
- lat
- lng
- country
- iso2_country_code
- population_estimate