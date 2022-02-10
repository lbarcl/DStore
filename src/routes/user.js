const SnowflakeId = require('../class/snowflake');
const router = require('express').Router();

const d = require('../class/discord')
const discord = new d(process.env.discord);
const userSchema = require('../database/userSchema')
const fileSchema = require('../database/fileSchema')
const { enc } = require('../util/base64')
const check = require('../util/authorize');


const snowflake = new SnowflakeId({
    offset: (2022 - 1970) * 31536000 * 1000
})

router.post('/new', async (req, res) => {
    try {
        const id = snowflake.generate()
        await userSchema({
            _id: id
        }).save()

        res.status(201).send({
            status: true,
            message: 'User created',
            data: {
                userId: id,
            }
        })
    } catch (err) {
        console.log(err)
        res.status(500).send(err);
    }
})

router.post('/:id/bucket/new', async (req, res) => {
    const userId = req.params.id
    try {
        if (!userId) {
            res.status(400).send({
                status: false,
                message: 'Missing user'
            });
        } else {
            const user = await userSchema.findById(userId)
            if (!user) {
                res.status(404).send({
                    status: false,
                    message: 'There is not a user with this id'
                })
            } else if (user.bucket) {
                res.status(406).send({
                    status: false,
                    message: 'This user already had a storage bucket'
                })
            } else {
                const response = await discord.createChannel(`${userId}-Bucket`)
                await userSchema.findByIdAndUpdate(userId, { bucket: response.id })
                let key = enc(`${process.env.secret}:${userId}:${response.id}:${Date.now()}`, "utf8")

                res.status(201).send({
                    status: true,
                    message: 'Bucket created',
                    data: {
                        bucketId: response.id,
                        bucketKey: key
                    }
                })
            }
        }
    } catch (err) {
        console.log(err)
        res.status(500).send(err);
    }
})

router.delete('/:id/bucket/:bucketId', check, async (req, res) => {
    const { id, bucketId } = req.params
    try {
        if (!id || !bucketId) {
            res.status(400).send({
                status: false,
                message: 'Missing user'
            });
        } else {
            const user = await userSchema.findById(id)
            if (!user) {
                res.status(404).send({
                    status: false,
                    message: 'There is not a user with this id'
                })
            } else if (user.bucket != bucketId) {
                res.status(403).send({
                    status: false,
                    message: 'This storage bucket is not yours'
                })
            } else {
                await discord.deleteChannel(bucketId)
                await fileSchema.deleteMany({ channelId: bucketId })
                await userSchema.findByIdAndUpdate(id, {bucket: ''})
                res.sendStatus(204)
            }
        }
    } catch (err) {
        console.log(err)
        res.status(500).send(err);
    }
})

module.exports = router