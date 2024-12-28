export namespace runtime {
	
	export class StoreProvider {
	
	
	    static createFrom(source: any = {}) {
	        return new StoreProvider(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	
	    }
	}
	export class FileSystem {
	
	
	    static createFrom(source: any = {}) {
	        return new FileSystem(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	
	    }
	}
	export class Browser {
	
	
	    static createFrom(source: any = {}) {
	        return new Browser(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	
	    }
	}
	export class Window {
	
	
	    static createFrom(source: any = {}) {
	        return new Window(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	
	    }
	}
	export class Dialog {
	
	
	    static createFrom(source: any = {}) {
	        return new Dialog(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	
	    }
	}
	export class Log {
	
	
	    static createFrom(source: any = {}) {
	        return new Log(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	
	    }
	}
	export class Events {
	
	
	    static createFrom(source: any = {}) {
	        return new Events(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	
	    }
	}
	export class Runtime {
	    // Go type: Events
	    Events?: any;
	    // Go type: Log
	    Log?: any;
	    // Go type: Dialog
	    Dialog?: any;
	    // Go type: Window
	    Window?: any;
	    // Go type: Browser
	    Browser?: any;
	    // Go type: FileSystem
	    FileSystem?: any;
	    // Go type: StoreProvider
	    Store?: any;
	
	    static createFrom(source: any = {}) {
	        return new Runtime(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.Events = this.convertValues(source["Events"], null);
	        this.Log = this.convertValues(source["Log"], null);
	        this.Dialog = this.convertValues(source["Dialog"], null);
	        this.Window = this.convertValues(source["Window"], null);
	        this.Browser = this.convertValues(source["Browser"], null);
	        this.FileSystem = this.convertValues(source["FileSystem"], null);
	        this.Store = this.convertValues(source["Store"], null);
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

