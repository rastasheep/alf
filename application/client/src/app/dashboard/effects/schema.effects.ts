import 'rxjs/add/operator/catch';
import 'rxjs/add/operator/map';
import 'rxjs/add/operator/switchMap';
import { Injectable } from '@angular/core';
import { Effect, Actions, toPayload } from '@ngrx/effects';
import { Action } from '@ngrx/store';
import { Observable } from 'rxjs/Observable';
import { empty } from 'rxjs/observable/empty';
import { of } from 'rxjs/observable/of';

import { SchemaService } from '../services/schema.service';
import { SchemaActions, LoadSuccessAction } from '../actions/schema.actions';
import { SchemaTable } from '../models/schema.model';

@Injectable()
export class SchemaEffects {
  @Effect()
  search$: Observable<Action> = this.actions$
    .ofType(SchemaActions.LOAD)
    .map(toPayload)
    .switchMap(_ => {
      return this.SchemaService
        .get()
        .map((schema: SchemaTable[]) => new LoadSuccessAction(schema))
        .catch(() => of(new LoadSuccessAction([])));
    });

  constructor(private actions$: Actions, private SchemaService: SchemaService) {}
}
