const { Schema, model } = require('mongoose');

const file = new Schema({
    _id: { type: String, required: true },
    bucket: String
})

module.exports = model('users', file);