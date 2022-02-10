const userSchema = require('../database/userSchema.js');
const { dec } = require('./base64');

async function check(req, res, next) {
    try {
        if (!req.header('authorization')) {
            res.status(400).send({
                status: false,
                message: 'No key sended'
            });
        } else {
            let key = dec(req.header('authorization'), "utf8")
            key = key.split(':')
            if (key[0] != process.env.secret) return res.sendStatus(401)

            const user = await userSchema.findById(key[1]);
            if (!user) {
                res.status(404).send({
                    status: false,
                    message: 'No user found'
                });
            } else if (key[2] != user.bucket) {
                res.status(403).send({
                    status: false,
                    message: 'This bucket is not yours'
                });
            } else {
                next()
            }   
        }
    } catch (err) {
        console.log(err)
        res.status(500).send(err);
    }
}

module.exports = check