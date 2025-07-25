export namespace models {
	
	export class MarketRating {
	    id: number;
	    overall_rating: number;
	    sector_ratings: Record<string, number>;
	    created_at: time.Time;
	    updated_at: time.Time;
	
	    static createFrom(source: any = {}) {
	        return new MarketRating(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.overall_rating = source["overall_rating"];
	        this.sector_ratings = source["sector_ratings"];
	        this.created_at = this.convertValues(source["created_at"], time.Time);
	        this.updated_at = this.convertValues(source["updated_at"], time.Time);
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
	export class MarketRatingRequest {
	    overall_rating: number;
	    sector_ratings: Record<string, number>;
	
	    static createFrom(source: any = {}) {
	        return new MarketRatingRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.overall_rating = source["overall_rating"];
	        this.sector_ratings = source["sector_ratings"];
	    }
	}
	export class OptionsTrade {
	    id: number;
	    ticker: string;
	    sector: string;
	    strategy_type: string;
	    entry_date: time.Time;
	    expiration_date: time.Time;
	    target_price?: number;
	    stop_loss?: number;
	    status: string;
	    notes: string;
	    created_at: time.Time;
	    updated_at: time.Time;
	
	    static createFrom(source: any = {}) {
	        return new OptionsTrade(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.ticker = source["ticker"];
	        this.sector = source["sector"];
	        this.strategy_type = source["strategy_type"];
	        this.entry_date = this.convertValues(source["entry_date"], time.Time);
	        this.expiration_date = this.convertValues(source["expiration_date"], time.Time);
	        this.target_price = source["target_price"];
	        this.stop_loss = source["stop_loss"];
	        this.status = source["status"];
	        this.notes = source["notes"];
	        this.created_at = this.convertValues(source["created_at"], time.Time);
	        this.updated_at = this.convertValues(source["updated_at"], time.Time);
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
	export class StrategyType {
	    id: number;
	    name: string;
	    category: string;
	    description: string;
	    color_hex: string;
	
	    static createFrom(source: any = {}) {
	        return new StrategyType(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.category = source["category"];
	        this.description = source["description"];
	        this.color_hex = source["color_hex"];
	    }
	}
	export class TradeRequest {
	    ticker: string;
	    sector: string;
	    strategy_type: string;
	    entry_date: time.Time;
	    expiration_date: time.Time;
	    target_price?: number;
	    stop_loss?: number;
	    notes: string;
	
	    static createFrom(source: any = {}) {
	        return new TradeRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ticker = source["ticker"];
	        this.sector = source["sector"];
	        this.strategy_type = source["strategy_type"];
	        this.entry_date = this.convertValues(source["entry_date"], time.Time);
	        this.expiration_date = this.convertValues(source["expiration_date"], time.Time);
	        this.target_price = source["target_price"];
	        this.stop_loss = source["stop_loss"];
	        this.notes = source["notes"];
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

export namespace time {
	
	export class Time {
	
	
	    static createFrom(source: any = {}) {
	        return new Time(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	
	    }
	}

}

