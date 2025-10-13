package url

type URLGetter interface {
	Get(id string) (*URL, error)
}

type URLSaver interface {
	Save(url *URL) error
}

//go:generate mockgen -destination=urlreaderwriter_mock.go -package=url . URLReaderWriter
type URLReaderWriter interface {
	URLGetter
	URLSaver
}

type ComputeShortFunc func(original string) (string, error)

type URLService struct {
	repo         URLReaderWriter
	computeShort ComputeShortFunc
}

func NewURLService(repo URLReaderWriter, computeShort ComputeShortFunc) *URLService {
	return &URLService{
		repo:         repo,
		computeShort: computeShort,
	}
}

func (h *URLService) Shorten(original string) (*URL, error) {
	short, err := h.computeShort(original)
	if err != nil {
		return nil, err
	}
	url := NewURL(short, original)
	if err := h.repo.Save(url); err != nil {
		return nil, err
	}
	return url, nil
}

func (h *URLService) Retrieve(id string) (*URL, error) {
	url, err := h.repo.Get(id)
	if err != nil {
		return nil, err
	}
	return url, nil
}
