CREATE TABLE users (
    email character varying(30) NOT NULL UNIQUE,
    password text NOT NULL,
    token text UNIQUE
);

CREATE TABLE users_confirmation (
    email character varying(30) NOT NULL UNIQUE,
    left_tries_count numeric DEFAULT 5 CHECK (left_tries_count <= 5 AND left_tries_count >= 0),
    confirmed boolean DEFAULT false,
    confirmation_code character(6)
);

CREATE TABLE general_chat_messages (
    id UUID NOT NULL,
    email character varying(30) NOT NULL,
    message_type character varying(30) NOT NULL,
    message_content text,
  	created_at timestamp
);
