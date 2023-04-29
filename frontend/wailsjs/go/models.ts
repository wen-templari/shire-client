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
	export class Group {
	    id?: number;
	    users?: User[];
	    createdAt?: string;
	    updatedAt?: string;
	
	    static createFrom(source: any = {}) {
	        return new Group(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.users = this.convertValues(source["users"], User);
	        this.createdAt = source["createdAt"];
	        this.updatedAt = source["updatedAt"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
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

