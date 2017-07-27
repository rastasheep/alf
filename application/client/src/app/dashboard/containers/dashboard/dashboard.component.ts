import { Component, ChangeDetectionStrategy, OnInit, OnDestroy } from '@angular/core';
import { ReactiveFormsModule, FormGroup, FormControl } from '@angular/forms';
import { ActivatedRoute, Router } from '@angular/router';
import { Store } from '@ngrx/store';
import 'rxjs/add/operator/takeUntil';
import { Observable } from 'rxjs/Observable';
import { Subject } from 'rxjs/Subject';

import { SchemaTable } from '../../models/schema.model';

@Component({
  selector: 'app-dashboard',
  changeDetection: ChangeDetectionStrategy.OnPush,
  styleUrls: ['./dashboard.component.css'],
  templateUrl: './dashboard.component.html'
})
export class DashboardComponent implements OnInit, OnDestroy {
  private ngUnsubscribe: Subject<void> = new Subject<void>();
  schema$: Observable<SchemaTable[]>;
  activeSection: string;
  id: number;

  constructor(
    private route: ActivatedRoute,
    private router: Router,
    private store: Store<any> // TODO
  ) {
    this.schema$ = store.select(state => state.dashboard.schema.tables);
  }

  ngOnInit() {
    this.route.data
      .takeUntil(this.ngUnsubscribe)
      .subscribe(r => this.activeSection = r.section);

    this.route.params
      .takeUntil(this.ngUnsubscribe)
      .subscribe((p) => this.id = +p.id);
  }

  ngOnDestroy() {
    this.ngUnsubscribe.next();
    this.ngUnsubscribe.complete();
  }
}
