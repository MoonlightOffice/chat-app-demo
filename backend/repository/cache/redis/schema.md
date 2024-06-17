# Redis DB schema

| Key | Domain object (in JSON format) |
| - | - |
| us.{user_id}.{session_id} | UserSession |
| l.{user_id} | Login |
| ur.{user_id} | []UserRoom |
| r.{room_id} | Room |
| rp.{room_id} | []Participant |
| m.{room_id}.{message_id} | Message |
