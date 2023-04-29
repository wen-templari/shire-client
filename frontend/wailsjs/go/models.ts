export namespace model {
	
	export class User {
	    id?: number;
	    name?: string;
	    address?: string;
	    port?: number;
	    rpcPort?: number;
	    password?: string;
	    createdAt?: string;
	    updatedAt?: string;
	
	    static createFrom(source: any = {}) {
	        return new User(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.address = source["address"];
	        this.port = source["port"];
	        this.rpcPort = source["rpcPort"];
	        this.password = source["password"];
	        this.createdAt = source["createdAt"];
	        this.updatedAt = source["updatedAt"];
	    }
	}
	export class Message {
	    from: number;
	    to: number;
	    content: string;
	    groupId: number;
	    time: string;
	    createdAt?: string;
	    updatedAt?: string;
	
	    static createFrom(source: any = {}) {
	        return new Message(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.from = source["from"];
	        this.to = source["to"];
	        this.content = source["content"];
	        this.groupId = source["groupId"];
	        this.time = source["time"];
	        this.createdAt = source["createdAt"];
	        this.updatedAt = source["updatedAt"];
	    }
	}

}

