module.exports = class {
    #senderFunc;

    constructor(file, Max) {
        this.file = file;
        this.fileParts = [];
        this.Max = Max;
    }

    DivideBytes() {
        let partLength = this.file.size / this.Max

        if (partLength < 1) partLength = 1;
        else if (partLength > 1) {
            partLength += ((partLength % this.Max) != 0) ? 1 : 0
            partLength = Math.floor(partLength)
        } 
        
        let start = 0;
        for (let i = 0; i < partLength; i++) {
            this.fileParts.push(this.file.data.slice(start, start + this.Max))
            start += this.Max
        }
    }

    Sender(func) {
        this.#senderFunc = func;
    }

    Send() {
        if (this.#senderFunc == null) throw new Error("Sender function must be defined before using this")
        else if (this.fileParts == null || this.fileParts.length == 0) throw new Error("Bytes must be divided before using this")

        this.fileParts.forEach(this.#senderFunc);
    }
};