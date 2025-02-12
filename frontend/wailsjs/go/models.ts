export namespace models {
	
	export class Account {
	    username: string;
	    publicKey: string;
	
	    static createFrom(source: any = {}) {
	        return new Account(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.username = source["username"];
	        this.publicKey = source["publicKey"];
	    }
	}
	export class Block {
	    Index: number;
	    Timestamp: number;
	    PrevHash: string;
	    Hash: string;
	    BlockType: string;
	    Data: any;
	
	    static createFrom(source: any = {}) {
	        return new Block(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Index = source["Index"];
	        this.Timestamp = source["Timestamp"];
	        this.PrevHash = source["PrevHash"];
	        this.Hash = source["Hash"];
	        this.BlockType = source["BlockType"];
	        this.Data = source["Data"];
	    }
	}
	export class Message {
	    sender: string;
	    receiver: string;
	    message: string;
	    timestamp: string;
	
	    static createFrom(source: any = {}) {
	        return new Message(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.sender = source["sender"];
	        this.receiver = source["receiver"];
	        this.message = source["message"];
	        this.timestamp = source["timestamp"];
	    }
	}

}

