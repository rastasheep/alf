import { Component, ChangeDetectionStrategy, OnInit, OnDestroy, Input } from '@angular/core';
import { Store } from '@ngrx/store';

import { SchemaTable } from '../../models/schema.model';
import { LoadAction } from '../../actions/schema.actions';

@Component({
  changeDetection: ChangeDetectionStrategy.OnPush,
  selector: 'app-schema',
  styleUrls: ['./schema.component.css'],
  templateUrl: './schema.component.html'
})
export class SchemaComponent implements OnInit, OnDestroy {
  @Input() schema: SchemaTable[];
  public schemaFilter = '';

  constructor(private store: Store<any>)  { }

  ngOnInit() {
    this.store.dispatch(new LoadAction());
  }

  ngOnDestroy() {
    // TODO: remove observer
  }
}
