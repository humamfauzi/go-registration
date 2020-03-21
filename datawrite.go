package registration

type DataWrite interface {
	Delete()
	DeleteBulk()
	Update()
	UpdateBulk()
}
