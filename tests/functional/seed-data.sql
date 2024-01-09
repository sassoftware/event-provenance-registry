INSERT INTO event_receivers
VALUES
('receiver-a', 'e2e tests', 'tests.e2e.finish', '1.0.0', 'completion of e2e tests', '{}', 'dk4n7ds', current_timestamp),
('receiver-b', 'SAST scan', 'scan.sast', '1.3.0', 'SAST scan of source code', '{}', 'n0kxa5g', current_timestamp),
('receiver-c', 'SCA scan', 'scan.sca', '1.2.4', 'SCA scan of source code', '{}', 'f3ksl2h', current_timestamp),
('receiver-d', 'publish artifact', 'artifact.publish', '2.0.0', 'publish to artifactory', '{}', 'vci6fk9', current_timestamp),
('receiver-e', 'manager sign-off', 'signoff.complete', '1.0.0', 'manual sign-off of artifact', '{"type": "object","properties": {"employee_id": {"type": "string"}}}', 're5u2al', current_timestamp);

INSERT INTO event_receiver_groups
VALUES
('group-a', 'deploy to prod', 'deploy.prod', '1.1.1', 'deploy to production when artifact vetted', true, current_timestamp, current_timestamp);

INSERT INTO event_receiver_group_to_event_receivers (event_receiver_id, event_receiver_group_id)
VALUES
('receiver-b', 'group-a'),
('receiver-c', 'group-a'),
('receiver-e', 'group-a');

INSERT INTO events
VALUES
('event-a', 'microservice-a', '1.2.3', '20201111.1605122728788', 'x64-oci-linux-2', 'docker', 'sast scan of microservice-a', '{}', true, current_timestamp, 'receiver-b'),
('event-b', 'microservice-a', '1.2.3', '20201111.1605122728788', 'x64-oci-linux-2', 'docker', 'sca scan of microservice-a', '{}', true, current_timestamp, 'receiver-c'),
('event-c', 'microservice-b', '2.7.5', '20231202.3938158958421', 'x64-oci-linux-2', 'docker', 'publish microservice-b to artifactory', '{}', true, current_timestamp, 'receiver-d');
