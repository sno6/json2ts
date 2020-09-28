export class BaseClass {
	 @IsArray() 
	public benchmarks!: Benchmarks[];
}

export class Benchmarks {
	 @IsString() 
	public benchmark_id!: string;

	public skill!: Skill;

	 @IsArray() 
	public metrics!: Metrics[];
}

export class Skill {
	 @IsNumber() 
	public rank!: number;

	 @IsString() 
	public category!: string;
}

export class Metrics {
	 @IsString() 
	public metric_id!: string;

	 @IsString() 
	public name!: string;

	 @IsNumber() 
	public value!: number;

	 @IsString() 
	public summary_field_mapping!: string;

	 @IsString() 
	public metric_id!: string;

	 @IsString() 
	public name!: string;

	 @IsNumber() 
	public value!: number;

	 @IsString() 
	public summary_field_mapping!: string;
}

