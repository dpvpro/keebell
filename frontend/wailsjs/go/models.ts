export namespace main {
	
	export class KDBXEntry {
	    UUID: string;
	    Title: string;
	    UserName: string;
	    Password: string;
	    URL: string;
	    Notes: string;
	    Icon: string;
	    Tags: string[];
	    CreatedTime: number;
	    ModifiedTime: number;
	    ExpiryTime?: number;
	    CustomFields: Record<string, string>;
	
	    static createFrom(source: any = {}) {
	        return new KDBXEntry(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.UUID = source["UUID"];
	        this.Title = source["Title"];
	        this.UserName = source["UserName"];
	        this.Password = source["Password"];
	        this.URL = source["URL"];
	        this.Notes = source["Notes"];
	        this.Icon = source["Icon"];
	        this.Tags = source["Tags"];
	        this.CreatedTime = source["CreatedTime"];
	        this.ModifiedTime = source["ModifiedTime"];
	        this.ExpiryTime = source["ExpiryTime"];
	        this.CustomFields = source["CustomFields"];
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

