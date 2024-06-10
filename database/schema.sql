CREATE DATABASE IF NOT EXISTS app;
USE app;

CREATE TABLE IF NOT EXISTS users (
  user_id VARCHAR(20) NOT NULL,
  created_at BIGINT UNSIGNED NOT NULL,

  PRIMARY KEY (user_id)
);

CREATE TABLE IF NOT EXISTS user_sessions (
  user_id VARCHAR(20) NOT NULL,
  session_id VARCHAR(20) NOT NULL,
  session_code TEXT NOT NULL,
  created_at BIGINT UNSIGNED NOT NULL,

  PRIMARY KEY (user_id, session_id),
  FOREIGN KEY (user_id) REFERENCES users (user_id)
);

CREATE TABLE IF NOT EXISTS logins (
  user_id VARCHAR(20) NOT NULL,
  login_code TEXT NOT NULL,
  created_at BIGINT UNSIGNED NOT NULL,

  PRIMARY KEY (user_id),
  FOREIGN KEY (user_id) REFERENCES users (user_id)
);

CREATE TABLE IF NOT EXISTS rooms (
  room_id VARCHAR(50) NOT NULL,
  created_at BIGINT UNSIGNED NOT NULL,
  name TEXT NOT NULL,
  last_message_id TEXT NOT NULL,

  PRIMARY KEY (room_id)
);

CREATE TABLE IF NOT EXISTS user_rooms (
  user_id VARCHAR(20) NOT NULL,
  room_id VARCHAR(50) NOT NULL,
  last_messaged_at BIGINT UNSIGNED NOT NULL,
  unread BIGINT UNSIGNED NOT NULL,

  PRIMARY KEY (user_id, room_id),
  FOREIGN KEY (user_id) REFERENCES users (user_id),
  FOREIGN KEY (room_id) REFERENCES rooms (room_id)
);

CREATE TABLE IF NOT EXISTS participants (
  room_id VARCHAR(50) NOT NULL,
  user_id VARCHAR(20) NOT NULL,
  joined_at BIGINT UNSIGNED NOT NULL,

  PRIMARY KEY (room_id, user_id),
  FOREIGN KEY (user_id) REFERENCES users (user_id),
  FOREIGN KEY (room_id) REFERENCES rooms (room_id)
);

CREATE TABLE IF NOT EXISTS messages (
  room_id VARCHAR(50) NOT NULL,
  message_id VARCHAR(50) NOT NULL,
  user_id VARCHAR(20) NOT NULL,
  content TEXT NOT NULL,
  sent_at BIGINT UNSIGNED NOT NULL,
  previous_message_id TEXT NOT NULL,

  PRIMARY KEY (room_id, message_id),
  FOREIGN KEY (room_id) REFERENCES rooms (room_id)
);
