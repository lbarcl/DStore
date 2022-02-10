const express = require('express');
const bodyParser = require('body-parser');
const mongoose = require('mongoose');

const app = express();

const PORT = process.env.PORT || 3000;

if (process.env.NODE_ENV ) {
    require('dotenv').config()
    console.log('Started on development envirionmet')
}

app.use(require('cors')());
app.use(bodyParser.json());
app.use(bodyParser.urlencoded({extended: true}));
app.use(require('morgan')('dev'));

app.use('/buckets', require('./routes/bucket'))
app.use('/users', require('./routes/user'))

app.listen(PORT, async () => {
    console.log('D-menager started listening on port ' + PORT)
    await mongoose.connect(process.env.mongo)
});