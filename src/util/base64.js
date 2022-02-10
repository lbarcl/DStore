function enc(data, encoding) {
    let buf = Buffer.from(data, encoding);
    return buf.toString("base64")
}

function dec(b64, decoding) {
    let buf = Buffer.from(b64, "base64");
    return buf.toString(decoding)
}

module.exports = {
    enc,
    dec
}