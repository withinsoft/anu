# Writing Commands

A Command in Anu is a binary in `$ANU/commands/*`. When a user sends a message
that starts with the command prefix (ex: `;foo` to run command `foo`), the 
correlating binary in `$ANU/commands` is run. Anu passes information to commands
with environment variables and the message body sent as standard input. Anything
written to standard output will be sent to the channel in question as command 
output. Anything written to standard error will be logged to the server console.

The following table lists all environment variables Anu populates and gives 
example values:

| Name                              | Type         | Example Value        | Description                                                                  |
| :----                             | :----        | :-------------       | :-----------                                                                 |
| `DISCORD_CHANNEL_ID`              | Snowflake ID | `267103998060789760` | The ID of the discord channel this command is being run in.                  |
| `DISCORD_CHANNEL_NAME`            | String       | `#foo`               | The human-readable name of the discord channel this command is being run in. |
| `DISCORD_GUILD_ID`                | Snowflake ID | `267103998060789760` | The ID of the discord guild this command is being run in.                    |
| `DISCORD_GUILD_NAME`              | String       | `The Source`         | The human-readable name of the discord guild this command is being run in.   |
| `DISCORD_MESSAGE_ID`              | Snowflake ID | `429079414916251658` | The ID of the discord message being processed.                               |
| `DISCORD_MESSAGE_AUTHOR_ID`       | Snowflake ID | `72838115944828928`  | The ID of the author of the message being processed.                         |
| `DISCORD_MESSAGE_AUTHOR_USERNAME` | String       | `Cadey~#1337`        | The Discord username+discriminator for the author of this message.           |
| `DISCORD_MESSAGE_AUTHOR_NICK`     | String       | `The Smiling Cadey`  | The nickname of this user in this guild.                                     |
| `USER`                            | String       | `Anu`                | The username of this bot user.                                               |
| `HOME`                            | String       | `/srv/anu`           | The location of Anu's state/command directory.                               |
| `VERB`                            | String       | `env`                | The command verb that started this handler.                                  |
