package blog

import (
	"testing"
	"time"
)

func TestCanServe(t *testing.T) {
	cases := []struct {
		name  string
		bl    Blog
		isDev bool
		want  bool
	}{
		{
			name: "draft in Dev",
			bl: Blog{
				IsDraft: true,
			},
			isDev: true,
			want:  true,
		},
		{
			name: "draft in prod",
			bl: Blog{
				IsDraft: true,
			},
			isDev: false,
			want:  false,
		},
		{
			name: "published",
			bl: Blog{
				IsDraft: false,
			},
			isDev: true,
			want:  true,
		},
		{
			name: "published",
			bl: Blog{
				IsDraft: false,
			},
			isDev: false,
			want:  true,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := tc.bl.CanServe(tc.isDev); got != tc.want {
				t.Errorf("CanServe() = %v, want %v", got, tc.want)
			}
		})
	}
}

var testBlogs = Blogs{
	Blog{
		Title:   "blog1",
		IsDraft: false,
		Date:    time.Date(2001, 1, 1, 1, 1, 1, 0, time.UTC),
	},
	Blog{
		Title:   "blog2",
		IsDraft: false,
		Date:    time.Date(2022, 1, 1, 1, 1, 1, 0, time.UTC),
	},
	Blog{
		Title:   "blog3",
		IsDraft: false,
		Date:    time.Date(2003, 1, 1, 1, 1, 1, 0, time.UTC),
	},
	Blog{
		Title:   "blog4",
		IsDraft: true,
		Date:    time.Date(2000, 1, 1, 1, 1, 1, 0, time.UTC),
	},
	Blog{
		Title:   "blog5",
		IsDraft: true,
		Date:    time.Date(2011, 1, 1, 1, 1, 1, 0, time.UTC),
	},
}

func TestGetNum(t *testing.T) {
	cases := []struct {
		name    string
		bl      Blogs
		n       int
		isDev   bool
		wantLen int
	}{
		{
			name:    "get all 5 blogs",
			bl:      testBlogs,
			n:       5,
			isDev:   true,
			wantLen: 5,
		},
		{
			name:    "get 5 blogs (prod) - 2 are drafts",
			bl:      testBlogs,
			n:       5,
			isDev:   false,
			wantLen: 3,
		},
		{
			name:    "get 2 blogs in order",
			bl:      testBlogs,
			n:       2,
			isDev:   false,
			wantLen: 2,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.bl.GetNum(tc.n, tc.isDev)

			if len(got) != tc.wantLen {
				t.Errorf("GetNum() = %v, want %v", len(got), tc.wantLen)
			}

			if !isSorted(got) {
				t.Errorf("GetNum() = %v, want sorted", got)
			}
		})
	}
}

func isSorted(blogs Blogs) bool {
	for i := 1; i < len(blogs); i++ {
		if !blogs[i-1].Date.After(blogs[i].Date) {
			return false
		}
	}
	return true
}
