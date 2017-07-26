import { Component, OnInit, OnDestroy } from '@angular/core';
import { Store } from '@ngrx/store';

import { LoadAction } from '../../actions/schema.actions';

@Component({
  selector: 'app-schema',
  styleUrls: ['./schema.component.css'],
  templateUrl: './schema.component.html'
})
export class SchemaComponent implements OnInit, OnDestroy {
  constructor(private store: Store<any>)  { }

  ngOnInit() {
    this.store.dispatch(new LoadAction());
  }

  ngOnDestroy() {
  }
}
