select erg.*, array_to_json(array(select event_receiver_id from event_receiver_group_to_event_receivers where event_receiver_group_id=erg.id)) as event_receiver_ids
from event_receiver_groups as erg
where id in
(
	select event_receiver_group_id
	from event_receiver_group_to_event_receivers
	where id not in
	(
		select id
		from
		(
			select gl.event_receiver_id
			, gl.event_receiver_group_id
			, CASE WHEN rec.success is NULL then false else rec.success END AS pass
			FROM (
				SELECT distinct sg.event_receiver_id, sg.event_receiver_group_id
				FROM event_receiver_group_to_event_receivers as sg
				WHERE event_receiver_group_id IN
				(
					SELECT event_receiver_group_id
					FROM event_receiver_group_to_event_receivers
					where event_receiver_id = @event_receiver_id
				)
			) as gl
			LEFT JOIN
			(
				SELECT DISTINCT ON (event_receiver_id) id, event_receiver_id, success, max(created_at) as ts
				FROM events
				where name = @name
				AND version = @version
				AND release = @release
				AND platform_id = @platform_id
				AND package = @package
				GROUP BY id
				ORDER BY event_receiver_id, ts DESC
			) as rec
			ON rec.event_receiver_id = gl.event_receiver_id
		) as false_gates
		where pass = false
	)
	and event_receiver_id=@event_receiver_id
)
and enabled = true
