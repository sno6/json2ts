export class BaseClass {
	@IsDefined()
	@IsArray()
	@Type(() => Data)
	@ValidateNested({ each: true })
	public data!: Data[];
}

export class Data {
	@IsDefined()
	@IsString()
	public status!: string;

	@IsDefined()
	@IsBoolean()
	public some-flag!: boolean;

	@IsDefined()
	@IsArray()
	@Type(() => Array)
	@ValidateNested({ each: true })
	public array!: Array[];
}

export class Array {
	@IsDefined()
	@IsString()
	public ofData!: string;
}

