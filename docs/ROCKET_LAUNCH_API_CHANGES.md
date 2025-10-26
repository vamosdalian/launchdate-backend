# Rocket Launch API Changes

## Overview
This document describes the changes made to the Rocket Launch API to align with the [RocketLaunch.Live API](https://www.rocketlaunch.live/api) structure.

## Changes Summary

### 1. RocketLaunch Model Extensions

The `RocketLaunch` model has been significantly extended to match the RocketLaunch.Live API structure:

#### New Fields

**Identifiers & Metadata:**
- `cospar_id` (string): COSPAR ID in format YYYY-nnn
- `sort_date` (string): Date for sorting purposes
- `slug` (string): URL-friendly identifier
- `modified` (timestamp): Last modification time

**Launch Window Fields:**
- `window_open` (timestamp): Launch window opening time
- `t0` (timestamp): Target launch time (T-0)
- `window_close` (timestamp): Launch window closing time
- `date_str` (string): Human-readable date string

**Provider Information:**
- `provider_id` (int64): Reference to the launch service provider (company)
- `provider` (nested object): Full provider details with id, name, and slug

**Vehicle Information:**
- `vehicle` (nested object): Rocket/vehicle details with id, name, company_id, and slug

**Pad & Location:**
- `pad` (nested object): Launch pad details with nested location information
  - `location`: Contains id, name, state, statename, country, and slug

**Mission Details:**
- `mission_description` (text): Description of the mission
- `launch_description` (text): Description of the launch
- `missions` (array): Array of mission objects with id, name, and description

**Weather Information:**
- `weather_summary` (text): Weather summary
- `weather_temp` (float32): Temperature
- `weather_condition` (string): Weather condition
- `weather_wind_mph` (float32): Wind speed in MPH
- `weather_icon` (string): Weather icon identifier
- `weather_updated` (timestamp): Last weather update time

**Additional Metadata:**
- `tags` (array): Array of tag objects with id and text
- `quicktext` (text): Quick summary text
- `suborbital` (boolean): Whether the launch is suborbital

### 2. Database Schema Changes

**Migration:** `003_update_rocket_launch_schema.up.sql`

#### Modified Tables

**rocket_launches table additions:**
- All new fields mentioned above as columns
- New indexes for performance:
  - `idx_rocket_launches_cospar_id`
  - `idx_rocket_launches_provider_id`
  - `idx_rocket_launches_slug`
  - `idx_rocket_launches_sort_date`
  - `idx_rocket_launches_t0`
  - `idx_rocket_launches_window_open`
  - `idx_rocket_launches_modified`

#### New Tables

**rocket_launch_missions:**
- `id` (bigserial): Primary key
- `rocket_launch_id` (bigint): Foreign key to rocket_launches
- `name` (varchar): Mission name
- `description` (text): Mission description
- Timestamps: `created_at`, `updated_at`

**rocket_launch_tags:**
- `id` (bigserial): Primary key
- `rocket_launch_id` (bigint): Foreign key to rocket_launches
- `text` (varchar): Tag text
- Timestamp: `created_at`

### 3. API Changes

#### Request/Response Models

**CreateRocketLaunchRequest** now includes all new fields:
```json
{
  "cospar_id": "2024-001A",
  "sort_date": "2024-01-15",
  "name": "Falcon 9 | Starlink",
  "provider_id": 1,
  "rocket_id": 1,
  "launch_base_id": 1,
  "mission_description": "Launch of Starlink satellites",
  "launch_description": "Falcon 9 launch from Cape Canaveral",
  "win_open": "2024-01-15T14:30:00Z",
  "t0": "2024-01-15T14:45:00Z",
  "win_close": "2024-01-15T15:00:00Z",
  "date_str": "Jan 15, 2024 14:45 UTC",
  "slug": "falcon-9-starlink-2024-001",
  "weather_summary": "Favorable conditions",
  "weather_temp": 72.5,
  "weather_condition": "Clear",
  "weather_wind_mph": 15.0,
  "weather_icon": "clear-day",
  "weather_updated": "2024-01-15T12:00:00Z",
  "quicktext": "SpaceX Falcon 9 launching Starlink satellites",
  "suborbital": false,
  "modified": "2024-01-15T13:00:00Z",
  "status": "scheduled"
}
```

**RocketLaunch Response** includes nested objects:
```json
{
  "id": 1,
  "cospar_id": "2024-001A",
  "name": "Falcon 9 | Starlink",
  "provider": {
    "id": 1,
    "name": "SpaceX",
    "slug": "spacex"
  },
  "vehicle": {
    "id": 1,
    "name": "Falcon 9 Block 5",
    "company_id": 1,
    "slug": "falcon-9-block-5"
  },
  "pad": {
    "id": 1,
    "name": "LC-39A",
    "location": {
      "id": 1,
      "name": "Kennedy Space Center",
      "state": "FL",
      "statename": "Florida",
      "country": "USA",
      "slug": "kennedy-space-center"
    }
  },
  "missions": [
    {
      "id": 1,
      "name": "Starlink",
      "description": "Deployment of Starlink satellites"
    }
  ],
  "tags": [
    {
      "id": 1,
      "text": "satellite"
    }
  ],
  "win_open": "2024-01-15T14:30:00Z",
  "t0": "2024-01-15T14:45:00Z",
  "win_close": "2024-01-15T15:00:00Z",
  "weather_summary": "Favorable conditions",
  "weather_temp": 72.5,
  "status": "scheduled",
  ...
}
```

#### Endpoints

All existing endpoints remain functional:
- `GET /api/v1/rocket-launches` - List rocket launches
- `POST /api/v1/rocket-launches` - Create a rocket launch
- `GET /api/v1/rocket-launches/{id}` - Get rocket launch details
- `PUT /api/v1/rocket-launches/{id}` - Update a rocket launch
- `DELETE /api/v1/rocket-launches/{id}` - Delete a rocket launch

The list endpoint now orders by `t0`, falling back to `window_open`, then `created_at`.

### 4. Repository Layer Changes

**RocketLaunchRepository** enhancements:
- Updated `Create()` to insert all new fields
- Updated `GetByID()` to load related entities (provider, vehicle, pad, missions, tags)
- Updated `List()` to load related entities for all results
- New `loadRelatedEntities()` helper method for nested data
- Improved ordering: `COALESCE(t0, window_open, created_at)`

### 5. Service Layer Changes

**RocketLaunchService** updates:
- `CreateRocketLaunch()` maps all new request fields to the model
- All cache invalidation logic remains intact

## Migration Instructions

To apply the migration:

```bash
# Using migrate CLI
migrate -path migrations -database "postgres://user:pass@localhost:5432/launchdate?sslmode=disable" up
```

To rollback:

```bash
migrate -path migrations -database "postgres://user:pass@localhost:5432/launchdate?sslmode=disable" down 1
```

## Frontend Adaptation Guide

When updating the frontend to use the new API structure:

1. **Update TypeScript/JavaScript interfaces** to include all new fields
2. **Handle nested objects**: Provider, Vehicle, Pad with Location
3. **Display mission array** instead of single mission description
4. **Show weather information** if available
5. **Use `t0` or `window_open`** for primary display time
6. **Display tags** as badges or chips
7. **Show launch window** (open to close time)
8. **Handle `quicktext`** for quick summaries

### Example Frontend Interface (TypeScript)

```typescript
interface RocketLaunch {
  id: number;
  cospar_id?: string;
  sort_date?: string;
  name: string;
  provider?: {
    id: number;
    name: string;
    slug: string;
  };
  vehicle?: {
    id: number;
    name: string;
    company_id?: number;
    slug: string;
  };
  pad?: {
    id: number;
    name: string;
    location?: {
      id: number;
      name: string;
      state: string;
      statename: string;
      country: string;
      slug: string;
    };
  };
  missions?: Array<{
    id: number;
    name: string;
    description: string;
  }>;
  mission_description?: string;
  launch_description?: string;
  win_open?: string;
  t0?: string;
  win_close?: string;
  date_str: string;
  tags?: Array<{
    id: number;
    text: string;
  }>;
  slug?: string;
  weather_summary?: string;
  weather_temp?: number;
  weather_condition?: string;
  weather_wind_mph?: number;
  weather_icon?: string;
  weather_updated?: string;
  quicktext?: string;
  suborbital: boolean;
  modified?: string;
  status: string;
  created_at: string;
  updated_at: string;
}
```

## Backward Compatibility

- All existing fields remain functional
- The `status` field still uses the same values: "scheduled", "successful", "failed", "cancelled"
- Old clients can continue to use basic fields; new fields are optional

## Future Enhancements

Consider adding query parameters to the list endpoint:
- `after_date` - Filter launches after a specific date
- `before_date` - Filter launches before a specific date
- `location_id` - Filter by launch location
- `pad_id` - Filter by launch pad
- `provider_id` - Filter by launch provider

These can be implemented by extending the repository `List()` method.

## Testing

All existing tests continue to pass. The changes are additive and don't break existing functionality.

To test the new fields:
1. Create a rocket launch with the new fields
2. Retrieve it and verify all nested objects are populated
3. List launches and verify sorting by t0/window_open
4. Update a launch with new fields

## Questions or Issues?

If you have questions about these changes or encounter issues when adapting the frontend, please open an issue in the repository.
