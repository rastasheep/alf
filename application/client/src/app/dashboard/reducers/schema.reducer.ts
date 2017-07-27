import { SchemaActions, SchemaAction } from '../actions/schema.actions';
import { SchemaTable } from '../models/schema.model';

export interface State {
  tables: SchemaTable[],
  openedTable: string,
}

export const initialState: State = {
  tables: [],
  openedTable: null,
};

export function SchemaReducer(state = initialState, action: SchemaAction): State {
  switch (action.type) {
    case SchemaActions.SELECT_TABLE: {
      return Object.assign({}, state, action.payload)
    }
    case SchemaActions.LOAD_SUCCESS: {
      return Object.assign({}, state, action.payload)
    }

    default: {
      return state;
    }
  }
}
