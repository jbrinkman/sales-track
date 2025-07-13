export namespace main {
	
	export class DatabaseHealth {
	    connected: boolean;
	    error?: string;
	
	    static createFrom(source: any = {}) {
	        return new DatabaseHealth(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.connected = source["connected"];
	        this.error = source["error"];
	    }
	}
	export class ImportError {
	    record: models.CreateSalesRecordRequest;
	    error: string;
	
	    static createFrom(source: any = {}) {
	        return new ImportError(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.record = this.convertValues(source["record"], models.CreateSalesRecordRequest);
	        this.error = source["error"];
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
	export class ImportOptions {
	    use_consignable_format: boolean;
	    custom_column_mapping?: string[];
	    strict_mode: boolean;
	    use_batch_import: boolean;
	
	    static createFrom(source: any = {}) {
	        return new ImportOptions(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.use_consignable_format = source["use_consignable_format"];
	        this.custom_column_mapping = source["custom_column_mapping"];
	        this.strict_mode = source["strict_mode"];
	        this.use_batch_import = source["use_batch_import"];
	    }
	}
	export class ImportResult {
	    success: boolean;
	    total_rows: number;
	    parsed_rows: number;
	    imported_rows: number;
	    error_message?: string;
	    parse_errors?: parser.ParseError[];
	    import_errors?: ImportError[];
	    processing_time: number;
	    imported_records?: models.SalesRecord[];
	    column_mapping: Record<string, number>;
	    data_types_detected: Record<string, string>;
	
	    static createFrom(source: any = {}) {
	        return new ImportResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.total_rows = source["total_rows"];
	        this.parsed_rows = source["parsed_rows"];
	        this.imported_rows = source["imported_rows"];
	        this.error_message = source["error_message"];
	        this.parse_errors = this.convertValues(source["parse_errors"], parser.ParseError);
	        this.import_errors = this.convertValues(source["import_errors"], ImportError);
	        this.processing_time = source["processing_time"];
	        this.imported_records = this.convertValues(source["imported_records"], models.SalesRecord);
	        this.column_mapping = source["column_mapping"];
	        this.data_types_detected = source["data_types_detected"];
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
	export class ImportStatistics {
	    total_records: number;
	    recent_records: number;
	    total_sales: number;
	    average_price: number;
	
	    static createFrom(source: any = {}) {
	        return new ImportStatistics(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.total_records = source["total_records"];
	        this.recent_records = source["recent_records"];
	        this.total_sales = source["total_sales"];
	        this.average_price = source["average_price"];
	    }
	}
	export class ValidationResult {
	    valid: boolean;
	    total_rows: number;
	    valid_rows: number;
	    invalid_rows: number;
	    error_message?: string;
	    errors?: parser.ParseError[];
	    warnings?: parser.ParseWarning[];
	    column_mapping: Record<string, number>;
	    data_types_detected: Record<string, string>;
	    processing_time: number;
	
	    static createFrom(source: any = {}) {
	        return new ValidationResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.valid = source["valid"];
	        this.total_rows = source["total_rows"];
	        this.valid_rows = source["valid_rows"];
	        this.invalid_rows = source["invalid_rows"];
	        this.error_message = source["error_message"];
	        this.errors = this.convertValues(source["errors"], parser.ParseError);
	        this.warnings = this.convertValues(source["warnings"], parser.ParseWarning);
	        this.column_mapping = source["column_mapping"];
	        this.data_types_detected = source["data_types_detected"];
	        this.processing_time = source["processing_time"];
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

export namespace models {
	
	export class CreateSalesRecordRequest {
	    store: string;
	    vendor: string;
	    date: string;
	    description: string;
	    sale_price: number;
	    commission: number;
	    remaining: number;
	
	    static createFrom(source: any = {}) {
	        return new CreateSalesRecordRequest(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.store = source["store"];
	        this.vendor = source["vendor"];
	        this.date = source["date"];
	        this.description = source["description"];
	        this.sale_price = source["sale_price"];
	        this.commission = source["commission"];
	        this.remaining = source["remaining"];
	    }
	}
	export class SalesRecord {
	    id: number;
	    store: string;
	    vendor: string;
	    // Go type: time
	    date: any;
	    description: string;
	    sale_price: number;
	    commission: number;
	    remaining: number;
	    // Go type: time
	    created_at: any;
	    // Go type: time
	    updated_at: any;
	
	    static createFrom(source: any = {}) {
	        return new SalesRecord(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.store = source["store"];
	        this.vendor = source["vendor"];
	        this.date = this.convertValues(source["date"], null);
	        this.description = source["description"];
	        this.sale_price = source["sale_price"];
	        this.commission = source["commission"];
	        this.remaining = source["remaining"];
	        this.created_at = this.convertValues(source["created_at"], null);
	        this.updated_at = this.convertValues(source["updated_at"], null);
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

export namespace parser {
	
	export class ParseError {
	    row: number;
	    column?: string;
	    message: string;
	    value?: string;
	
	    static createFrom(source: any = {}) {
	        return new ParseError(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.row = source["row"];
	        this.column = source["column"];
	        this.message = source["message"];
	        this.value = source["value"];
	    }
	}
	export class ParseWarning {
	    row: number;
	    column?: string;
	    message: string;
	    value?: string;
	
	    static createFrom(source: any = {}) {
	        return new ParseWarning(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.row = source["row"];
	        this.column = source["column"];
	        this.message = source["message"];
	        this.value = source["value"];
	    }
	}

}

