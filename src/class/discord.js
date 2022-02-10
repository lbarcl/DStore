const fetch = require('node-fetch')

class discord {
    constructor(key) {
        if (!key) throw new Error("No key provided")
        this.key = key
        this.baseUrl = "https://discord.com/api/v9/"
    }

    async sendMessage(channel, content) {
        let response = await fetch(`${this.baseUrl}channels/${channel}/messages`, {
            method: "POST",
            headers: {
                Authorization: 'Bot ' + this.key,
                "Content-Type": "application/json"
            },
            body: JSON.stringify({
                content,
                tts: false,
            })
        })

        response = await response.json()
        return response
    }

    async sendMultiPart(channel, form) {
        let response = await fetch(`${this.baseUrl}channels/${channel}/messages`, {
            method: 'POST',
            body: form,
            headers: { 
                Authorization: 'Bot ' + this.key
            }
        })

        response = await response.json()
        return response
    }

    async downloadAttachment(channel, attachment, filename) {
        const url = `https://cdn.discordapp.com/attachments/${channel}/${attachment}/${filename}`

        let response = await fetch(url, { method: 'GET' })
        if (!response.ok) {
            throw response.error()
        }
        response = await response.buffer()
        return response
    }

    async createChannel(channelName) {
        const response = await fetch(`${this.baseUrl}guilds/788450913136148600/channels`, {
            method: 'POST',
            headers: {
                Authorization: 'Bot ' + this.key,
                "Content-Type": "application/json"
            },
            body: JSON.stringify({
                name: channelName,
                type: 0,
                parent_id: '930204410850738176',
            })
        })

        const jsonRes = await response.json()
        return jsonRes
    }

    async deleteMessage(channel, message) {
        let response = await fetch(`${this.baseUrl}channels/${channel}/messages/${message}`, {
            method: 'DELETE',
            headers: {
                Authorization: 'Bot ' + this.key,
                "X-Audit-Log-Reason": 'Delete request by owner of the file'
            }
        })

        if (response.status != 204) {
            response = await response.json()
            return response
        } else {
            return 204
        }
    }

    async deleteChannel(channel) {
        let response = await fetch(`${this.baseUrl}channels/${channel}`, {
            method: 'DELETE',
            headers: {
                Authorization: 'Bot ' + this.key,
                "X-Audit-Log-Reason": 'Delete request by owner of the bucket'
            }
        })

        if (response.status != 204) {
            response = await response.json()
            return response
        } else {
            return 204
        }
    }
}

module.exports = discord