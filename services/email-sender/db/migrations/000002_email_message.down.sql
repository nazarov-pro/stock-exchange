ALTER TABLE "EMAIL"."EMAIL_MESSAGES" DROP CONSTRAINT "FK_EMAIL_MESSAGE_RECIPIENT_EMAIL";
DROP TABLE "EMAIL"."EMAIL_MESSAGE_RECIPIENTS";
DROP TABLE "EMAIL"."EMAIL_MESSAGES";
DROP SEQUENCE "EMAIL"."EMAIL_MESSAGES_ID_SEQ";