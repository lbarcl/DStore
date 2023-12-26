# DStore

DStore is a project developed in GoLang, designed for educational purposes only. It functions as a file storage system on Discord, utilizing the platform's storage space. The initial idea was to provide a potential alternative to services like Google Drive, with the belief that it could offer faster performance. However, the project is not actively supported or updated.

## Purpose
The primary purpose of DStore is to store files on Discord, essentially utilizing Discord's storage space. It was conceived as an experimental project to explore the possibilities of creating a decentralized file storage solution within the Discord ecosystem.

## Functionality
DStore operates through a set of commands, each serving a specific purpose:
- **help:** Displays a list of available commands.
- **upload:** Initiates the process of uploading a file to Discord. The upload process involves creating a new file structure, calculating necessary parts (messages) for upload, and handling potential disconnections due to rate limiting by Discord.
- **resume:** Allows users to resume uploading a file using the file ID, useful in case of disconnections during the upload process.
- **list:** Displays a list of files stored on Discord through DStore.
- **exit:** Closes the DStore application.

## Upload Mechanism
Under the hood, DStore employs a systematic approach to file uploads:
1. **File Structure:** Upon initiating an upload, a new file structure is created to track the progress and details of the upload.
2. **Parts Calculation:** The system calculates the number of parts (messages) needed for the upload. Discord's policy limits attachments to 25MB per message.
3. **Chunked Upload:** Instead of buffering the entire file to avoid stack overflow, DStore chunks the file and sends parts to Discord as attachments.
4. **Fail-Safe Mechanism:** Due to potential disconnections from Discord's server, a fail-safe mechanism allows users to resume the upload using the file ID.
5. **Storage:** Completed uploads are saved, including message and attachment IDs, in the parts folder of each file as JSON (1KB per part).

## Code Overview
The code for DStore is designed to be simple and self-explanatory, offering transparency into its functionality. The approach to file upload and handling potential issues is outlined to facilitate understanding. You need to change the configuration for Discord under the `/Discord/discord.go` file.

## Note
Prior to the GoLang implementation, there was a NodeJS server code that served a similar purpose. Additionally, there exists a sister repository designed to interact with the NodeJS server.

**Disclaimer:** This project is intended for educational purposes only and is not actively maintained or updated. Users are encouraged to use this project responsibly and in compliance with Discord's terms of service.