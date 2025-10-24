# Rocket Launch Sync Endpoint

## Overview

This document describes the new endpoint for syncing rocket launch data from the RocketLaunch.Live API.

## Endpoint

**POST** `/api/v1/rocket-launches/sync`

This endpoint fetches the latest 5 upcoming rocket launches from [RocketLaunch.Live API](https://fdo.rocketlaunch.live/json/launches/next/5) and saves them to the database.

## Request

No request body is required.

**Example:**
```bash
curl -X POST http://localhost:8080/api/v1/rocket-launches/sync
```

## Response

**Success (200 OK):**
```json
{
  "message": "rocket launches synced successfully",
  "count": 5
}
```

Where `count` is the number of launches that were successfully created or updated in the database.

**Error (500 Internal Server Error):**
```json
{
  "error": "failed to sync rocket launches"
}
```

## Behavior

1. **Fetch**: The endpoint fetches data from `https://fdo.rocketlaunch.live/json/launches/next/5`
2. **Parse**: The response is parsed into the internal data model
3. **Save**: For each launch:
   - Attempts to create a new record in the database
   - If creation fails (e.g., duplicate slug), attempts to update the existing record by slug
   - Increments the success counter if saved/updated successfully
4. **Cache**: Invalidates the rocket launches cache to ensure fresh data is returned
5. **Return**: Returns the count of successfully synced launches

## Data Mapping

The following fields are mapped from the external API to our database:

- `cospar_id` - COSPAR ID
- `sort_date` - Sorting date
- `name` - Launch name
- `mission_description` - Mission description
- `launch_description` - Launch description
- `win_open` - Window open time
- `t0` - Target launch time (T-0)
- `win_close` - Window close time
- `date_str` - Human-readable date string
- `slug` - URL-friendly identifier (used for duplicate detection)
- `weather_summary` - Weather summary
- `weather_temp` - Temperature
- `weather_condition` - Weather condition
- `weather_wind_mph` - Wind speed in MPH
- `weather_icon` - Weather icon identifier
- `weather_updated` - Weather update timestamp
- `quicktext` - Quick summary text
- `suborbital` - Suborbital flag
- `modified` - Last modification time
- `status` - Always set to "scheduled" for synced launches

## Important Notes

1. **External IDs**: Provider, vehicle, and launch base IDs from the external API are **not** mapped to our database because they use different ID systems. These relationships should be established manually or through a separate matching process.

2. **Duplicate Handling**: Duplicates are detected using the `slug` field. If a launch with the same slug already exists, it will be updated with the latest data.

3. **Status**: All synced launches are set to "scheduled" status by default.

4. **Cache Invalidation**: The endpoint automatically invalidates all rocket launch caches to ensure fresh data is returned in subsequent queries.

## Use Cases

- **Scheduled Sync**: Can be called periodically (e.g., via cron job) to keep the database up-to-date with the latest launch information
- **Manual Refresh**: Can be triggered manually when fresh launch data is needed
- **Initial Data Load**: Can be used to populate the database with upcoming launches

## Example Usage

### Using curl
```bash
curl -X POST http://localhost:8080/api/v1/rocket-launches/sync
```

### Using HTTPie
```bash
http POST http://localhost:8080/api/v1/rocket-launches/sync
```

### Using JavaScript/Fetch
```javascript
fetch('http://localhost:8080/api/v1/rocket-launches/sync', {
  method: 'POST',
})
  .then(response => response.json())
  .then(data => console.log(`Synced ${data.count} launches`))
  .catch(error => console.error('Error:', error));
```

## Future Enhancements

Potential improvements for this endpoint:

1. **Configurable Count**: Allow specifying how many launches to fetch (instead of hardcoded 5)
2. **Date Range Filtering**: Allow filtering by date range
3. **Provider Matching**: Automatically match external provider IDs to our database
4. **Batch Operations**: Use database transactions for atomic batch inserts/updates
5. **Webhook Support**: Trigger sync automatically when external API has updates
6. **Conflict Resolution**: More sophisticated handling of conflicts when updating existing launches
