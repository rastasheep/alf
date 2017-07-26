import { Action } from '@ngrx/store';
import { SchemaTable } from '../models/schema.model';

export const SchemaActions = {
  SELECT_TABLE: '[SCHEMA] SELECT_TABLE',
  LOAD_SUCCESS: '[SCHEMA] LOAD_SUCCESS',
  LOAD: '[SCHEMA] LOAD',
};

export class LoadAction implements Action {
  readonly type = SchemaActions.LOAD;
  payload: {}
}

export class LoadSuccessAction implements Action {
  readonly type = SchemaActions.LOAD_SUCCESS;
  payload: { tables: SchemaTable[] }

  constructor(tables: SchemaTable[]) {
    this.payload = { tables };
  }
}

export class SelectTableAction implements Action {
  readonly type = SchemaActions.SELECT_TABLE;
  payload: { openedTable: string }

  constructor(tableName: string) {
    this.payload = {
      openedTable: tableName
    };
  }
}

export type SchemaAction = LoadAction | LoadSuccessAction | SelectTableAction;

