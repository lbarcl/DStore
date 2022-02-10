const fileParser = require('express-fileupload');
const router = require('express').Router();
const formBody = require('form-data');

const d = require('../class/discord')
const discord = new d(process.env.discord);
const fileSchema = require('../database/fileSchema.js');
const file = require('../class/file.js');
const check = require('../util/authorize');

router.use(fileParser({
    createParentPath: true,
}))

router.post('/:bucketId/files/new', check, async (req, res) => {
    try {
        if (!req.files || !req.files.file) {
            //* No file uploaded
            res.status(400).send({
                status: false,
                message: 'No file uploaded'
            });
        }
        else if (!req.params.bucketId) {
            //* No channel ID provided
            res.status(400).send({
                status: false,
                message: 'No bucket ID provided'
            });
        } else {
            //* File implamentions
            const f = new file(req.files.file, 8388608) //2MB -> 2097152 | 8MB -> 8388608
            f.DivideBytes()

            //* Sending file chunks to discord
            const body = new formBody()
            let name = f.file.name.split('.')
            name.pop()
            name = name.join('.')

            for (let i = 0; i < f.fileParts.length; i++) {
                body.append(`file[${i}]`, f.fileParts[i], name + `-${i}.fbin`)
            } 

            const response = await discord.sendMultiPart(req.params.bucketId, body)
            
            //* Saving data to mongo db
            const data = {
                _id: response.id,
                fileType: f.file.mimetype,
                fileName: f.file.name,
                fileSize: f.file.size,
                channelId: response.channel_id,
                fileChunks: [],
            }

            for (let at of response.attachments) {
                data.fileChunks.push({
                    _id: at.id,
                    fileName: at.filename,
                })
            }

            await fileSchema(data).save();

            //* Sending response
            res.status(201).send({
                status: true,
                message: 'File is uploaded',
                data: {
                    bucketId: data.channelId,
                    fileId: data._id,
                    fileName: data.fileName,
                    fileSize: data.fileSize,
                    fileType: data.fileType,
                    fileUrl: `http://localhost:3000/bucket/${data.channelId}/${data._id}/${data.fileName}`
                }
            });
        }
    } catch (err) {
        console.log(err)
        res.status(500).send(err);
    }
})

router.get('/:bucketId/files/:fileId', async (req, res) => {
    const {bucketId, fileId} = req.params
    try {
        if (!bucketId || !fileId) {
            res.status(400).send({
                status: false,
                message: 'Missing path parameter'
            });
        } else {
            const fileMeta = await fileSchema.findOne({ _id: fileId, channelId: bucketId })
            if (!fileMeta) {
                res.status(404).send({ status: false, message: 'Couldn\'t find any file' })
            } else {
                const chunks = []
                
                for (let i = 0; i < fileMeta.fileChunks.length; i++) {
                    try {
                        const fbin = await discord.downloadAttachment(bucketId, fileMeta.fileChunks[i]._id, fileMeta.fileChunks[i].fileName)
                        chunks.push(fbin)
                    } catch (error) {
                        console.log(error)
                        res.status(404).send({ status: false, message: 'Couldn\'t find any file' })
                        return
                    }
                }

                const data = (chunks.length > 1) ? Buffer.concat(chunks) : chunks[0]
                
                res.setHeader('content-type', fileMeta.fileType)
                res.setHeader('content-length', data.length)
                if (req.query.d) res.setHeader('Content-Disposition', `attachment; filename="${fileMeta.fileName}"`)
                
                const f = new file({ data, size: data.length }, 65475)
                f.DivideBytes()
                
                for (let i = 0; i < f.fileParts.length; i++) {
                    res.write(f.fileParts[i])
                }
                
                res.end()
            }
        }
    } catch (err) {
        console.log(err)
        res.status(500).send(err);
    }
})

router.get('/:bucketId/files/:fileId/meta', async (req, res) => {
    const {bucketId, fileId} = req.params
    try {
        if (!bucketId || !fileId) {
            res.status(400).send({
                status: false,
                message: 'Missing path parameter'
            });
        } else {
            const fileMeta = await fileSchema.findOne({ _id: fileId, channelId: bucketId})
            if (!fileMeta) {
                res.status(404).send({ status: false, message: 'Couldn\'t find any file' })
            } else {
                res.setHeader('Content-Type', 'application/json')
                res.send(JSON.stringify({
                    id: fileMeta._id,
                    fileType: fileMeta.fileType,
                    fileName: fileMeta.fileName,
                    fileSize: fileMeta.fileSize
                }))
            }
        }
    } catch (err) {
        console.log(err)
        res.status(500).send(err);
    }
})

router.get('/:bucketId/files', check, async (req, res) => {
    const {bucketId} = req.params
    try {
        if (!bucketId) {
            res.status(400).send({
                status: false,
                message: 'Missing path parameter'
            });
        } else {
            const files = await fileSchema.find({channelId: bucketId})
            if (!files || files.length == 0) {
                res.status(404).send({ status: false, message: 'Couldn\'t find any file' })
            } else {
                let metaFiles = files.map(fileMeta => {
                    return {
                        id: fileMeta._id,
                        fileType: fileMeta.fileType,
                        fileName: fileMeta.fileName,
                        fileSize: fileMeta.fileSize
                    }
                })
                res.setHeader('Content-Type', 'application/json')
                res.send(JSON.stringify(metaFiles))
            }
        }
    } catch (err) {
        console.log(err)
        res.status(500).send(err);
    }
})

router.delete('/:bucketId/files/:fileId', check, async (req, res) => {
    const { bucketId, fileId } = req.params
    try {
        if (!bucketId || !fileId) {
            res.status(400).send({
                status: false,
                message: 'Missing path parameter'
            });
        } else {
            const fileMeta = await fileSchema.findOne({ _id: fileId, channelId: bucketId})
            if (!fileMeta) {
                res.status(404).send({ status: false, message: 'Couldn\'t find any file' })
            } else {
                const response = await discord.deleteMessage(bucketId, fileId)
                if (response === 204) {
                    await fileSchema.findOneAndRemove({ _id: fileId, channelId: bucketId})
                    res.sendStatus(204)
                } else {
                    await fileSchema.findOneAndRemove({ _id: fileId, channelId: bucketId})
                    res.status(404).send({ status: false, message: 'Couldn\'t find any file on the storage'})
                }
            }
        }
    } catch (err) {
        console.log(err)
        res.status(500).send(err);
    }
})

module.exports = router