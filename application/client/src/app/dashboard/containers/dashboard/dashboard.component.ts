import { Component, OnInit, OnDestroy } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import 'rxjs/add/operator/takeUntil';
import { Subject } from 'rxjs/Subject';

@Component({
  selector: 'app-dashboard',
  styleUrls: ['./dashboard.component.css'],
  templateUrl: './dashboard.component.html'
})
export class DashboardComponent implements OnInit, OnDestroy {
  private ngUnsubscribe: Subject<void> = new Subject<void>();
  sub: any;
  activeSection: string;
  id: number;

  constructor(
    private route: ActivatedRoute,
    private router: Router,
  ) {
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
