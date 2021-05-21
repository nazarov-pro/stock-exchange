CREATE SEQUENCE "EMAIL"."EMAIL_MESSAGES_ID_SEQ" AS BIGINT INCREMENT BY 1;
CREATE TABLE "EMAIL"."EMAIL_MESSAGES" (
    "ID" BIGINT PRIMARY KEY,
    "SUBJECT" VARCHAR(127),
    "SENDER" VARCHAR(127),
    "CONTENT" TEXT NOT NULL,
    "STATUS" SMALLINT NOT NULL,
    "CREATED_DATE" BIGINT NOT NULL,
    "LAST_UPDATE_DATE" BIGINT
);
CREATE TABLE "EMAIL"."EMAIL_MESSAGE_RECIPIENTS" (
    "EMAIL_MSG_ID" BIGINT NOT NULL,
    "RECIPIENT" VARCHAR(127) NOT NULL,
    PRIMARY KEY ("EMAIL_MSG_ID", "RECIPIENT")
);
ALTER TABLE "EMAIL"."EMAIL_MESSAGE_RECIPIENTS"
ADD CONSTRAINT "FK_EMAIL_MESSAGE_RECIPIENT_EMAIL" FOREIGN KEY ("EMAIL_MSG_ID") REFERENCES "EMAIL"."EMAIL_MESSAGES"("ID");