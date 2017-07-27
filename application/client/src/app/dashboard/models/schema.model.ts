export interface SchemaTable {
  tableName: string;
  columns: {
    tableName: string;
    columnName: string;
    dataType: string;
  }[];
}
