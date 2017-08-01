import { Action } from '@ngrx/store';
import { Execution } from '../models/execution.model';

export const ExecutionActions = {
  LOAD_SINGLE_SUCCESS: '[EXECUTION] LOAD_SINGLE_SUCCESS',
  LOAD_SUCCESS: '[EXECUTION] LOAD_SUCCESS',
  LOAD: '[EXECUTION] LOAD',
};

export class LoadSingleSuccessAction implements Action {
  readonly type = ExecutionActions.LOAD_SINGLE_SUCCESS;
  payload: { execution: Execution }

  constructor(execution: Execution) {
    this.payload = { execution };
  }
}

export class LoadSuccessAction implements Action {
  readonly type = ExecutionActions.LOAD_SUCCESS;
  payload: { executions: Execution[] }

  constructor(executions: Execution[]) {
    this.payload = { executions };
  }
}

export class LoadAction implements Action {
  readonly type = ExecutionActions.LOAD;
  payload: {
    id: string,
    lastId: number,
  }

  constructor(id: string, lastId: number) {
    this.payload = { id, lastId };
  }
}

export type ExecutionAction = LoadAction;

