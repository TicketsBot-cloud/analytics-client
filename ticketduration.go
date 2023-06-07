package analytics

import (
	"context"
	"time"
)

func (c *Client) GetTicketDurationStats(context context.Context, guildId uint64) (TripleWindow, error) {
	query := `
SELECT
    avgMerge(all_time),
    avgOrNullMerge(monthly),
    avgOrNullMerge(weekly)
FROM analytics.ticket_duration
WHERE guild_id = ?
GROUP BY guild_id`

	rows, err := c.client.Query(context, query, guildId)
	if err != nil {
		return TripleWindow{}, err
	}

	if rows.Next() {
		// Values in seconds
		var allTime int64
		var monthly, weekly *int64
		if err := rows.Scan(&allTime, &monthly, &weekly); err != nil {
			return TripleWindow{}, err
		}

		return TripleWindow{
			AllTime: ptr(time.Duration(allTime) * time.Second),
			Monthly: mapNullableSecondsToDuration(monthly),
			Weekly:  mapNullableSecondsToDuration(weekly),
		}, nil
	} else {
		return blankTripleWindow(), nil
	}
}
