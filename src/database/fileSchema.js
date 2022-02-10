const { Schema, model } = require('mongoose');

const fileChunk = new Schema({
    _id: { type: String, required: true },
    fileName: { type: String, required: true }
})


const file = new Schema({
    _id: { type: String, required: true },
    fileType: { type: String, default: 'application/octet-stream' },
    fileName: { type: String, required: true },
    channelId: { type: String, required: true },
    fileChunks: { type: [fileChunk], required: true },
    fileSize: Number,
})

module.exports = model('files', file);