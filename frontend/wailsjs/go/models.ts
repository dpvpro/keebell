export namespace main {
	
	export class KDBXEntry {
	    uuid: string;
	    title: string;
	    userName: string;
	    password: string;
	    url: string;
	    notes: string;
	    icon: string;
	    tags: string[];
	    createdTime: number;
	    modifiedTime: number;
	    expiryTime?: number;
	    customFields?: Record<string, string>;
	
	    static createFrom(source: any = {}) {
	        return new KDBXEntry(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.uuid = source["uuid"];
	        this.title = source["title"];
	        this.userName = source["userName"];
	        this.password = source["password"];
	        this.url = source["url"];
	        this.notes = source["notes"];
	        this.icon = source["icon"];
	        this.tags = source["tags"];
	        this.createdTime = source["createdTime"];
	        this.modifiedTime = source["modifiedTime"];
	        this.expiryTime = source["expiryTime"];
	        this.customFields = source["customFields"];
	    }
	}
	export class KDBXGroup {
	    uuid: string;
	    name: string;
	    icon: string;
	    entries?: KDBXEntry[];
	    groups?: KDBXGroup[];
	
	    static createFrom(source: any = {}) {
	        return new KDBXGroup(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.uuid = source["uuid"];
	        this.name = source["name"];
	        this.icon = source["icon"];
	        this.entries = this.convertValues(source["entries"], KDBXEntry);
	        this.groups = this.convertValues(source["groups"], KDBXGroup);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
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

}

