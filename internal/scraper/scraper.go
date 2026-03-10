package scraper

import (
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"

	"utnianos-api/internal/models"

	"github.com/PuerkitoBio/goquery"
)

const BaseURL = "https://www.utnianos.com.ar/foro/"
const TimeLayout = "02-01-2006 15:04"

var DateTextRegex = regexp.MustCompile(`\d{2}-\d{2}-\d{4} \d{2}:\d{2}`)

type Scraper struct {
	client *http.Client
}

func New() *Scraper {
	return &Scraper{
		client: &http.Client{
			Timeout: 30 * time.Second,
			//CheckRedirect: func(req *http.Request, via []*http.Request) error {
			//	return http.ErrUseLastResponse
			//},
		},
	}
}

func (s *Scraper) GetForums() (*models.ForumListResponse, error) {
	doc, err := s.fetchDoc(BaseURL + "index.php")
	if err != nil {
		return nil, err
	}

	response := &models.ForumListResponse{
		Forums: []models.Forum{},
	}

	doc.Find("td > strong a[href*='foro-']").Each(func(_ int, sel *goquery.Selection) {
		forum, err := s.parseForumRow(sel)
		if err != nil {
			return
		} // Continua con el siguiente foro si hay un error
		response.Forums = append(response.Forums, forum)
	})

	return response, nil
}

func (s *Scraper) parseForumRow(sel *goquery.Selection) (models.Forum, error) {
	forum := models.Forum{}

	href, exists := sel.Attr("href")
	if !exists {
		return forum, fmt.Errorf("forum link not found")
	}

	href, _ = url.QueryUnescape(href) // Manejo de caracteres especiales en la URL

	forum.URL = href
	forum.Name = strings.TrimSpace(sel.Text())

	parts := strings.Split(href, "foro-")
	if len(parts) > 1 {
		forum.Slug = parts[1]
	}

	// funcion aux para parsear numeros ("1.234" -> 1234)
	parseNumber := func(s string) (int, error) {
		return strconv.Atoi(strings.ReplaceAll(strings.TrimSpace(s), ".", ""))
	}

	// Como si fuera un puntero, al <td> del enlace, para recorrer los <td> siguientes
	// Me "muevo" entre los campos de la tabla
	td := sel.Closest("td").Next() // numero de temas
	if topic_count, err := parseNumber(td.Text()); err == nil {
		forum.TopicCount = topic_count
	}

	td = td.Next() // numero de mensajes
	if message_count, err := parseNumber(td.Text()); err == nil {
		forum.MessageCount = message_count
	}

	span := td.Next().Children().First() // Ultimo post
	dateText := DateTextRegex.FindString(span.Text())
	fmt.Println(dateText)
	lastPostDate, _ := time.Parse(TimeLayout, dateText)
	href, _ = span.Children().First().Attr("href")
	href, _ = url.QueryUnescape(href)

	author := span.Children().Last()
	authorURL, _ := author.Attr("href")
	authorURL, _ = url.QueryUnescape(authorURL)

	forum.LastPost = &models.ForumLastPost{
		URL:       BaseURL + href,
		Date:      &lastPostDate,
		Author:    author.Text(),
		AuthorURL: authorURL,
	}

	return forum, nil
}

func (s *Scraper) GetForum(slug string, filters map[string][]string) (*models.SearchResult, error) {
	if slug == "" {
		return nil, fmt.Errorf("forum slug is required")
	}

	forumURL := BaseURL + "foro-" + slug + "?"
	q := url.Values{}
	result := &models.SearchResult{
		Topics: []models.Topic{}, // Lo rellena parseForumPage
		Page:   1,
	}

	for key, values := range filters {
		for _, v := range values {
			q.Add(key, v)
		}
	}

	if len(q) > 0 {
		queryParams, _ := url.QueryUnescape(q.Encode())
		forumURL += queryParams
	}

	fmt.Println("Fetching forum URL: ", forumURL)
	doc, err := s.fetchDoc(forumURL)
	if err != nil {
		return nil, err
	}

	// La website de utnianos devuelve la página 1, si page=foo no es válido
	if page, err := strconv.Atoi(q.Get("page")); err == nil && page > 1 {
		result.Page = page
	}

	return s.parseForumPage(doc, result)
}

func (s *Scraper) parseForumPage(doc *goquery.Document, result *models.SearchResult) (*models.SearchResult, error) {
	// Paginas totales
	selPageTotal := doc.Find(".pages")
	totalPageRegex := regexp.MustCompile(`\(([^)]+)\)`)
	totalPageStr := totalPageRegex.FindStringSubmatch(selPageTotal.First().Text())

	if totalPage, err := strconv.Atoi(totalPageStr[1]); err == nil {
		result.TotalPages = totalPage
	}

	// Recupero y parseo los temas encontrados de la página.
	doc.Find("table.tborder tr").Each(func(_ int, sel *goquery.Selection) {
		omitClasses := []string{"thead", "tcat", "tfoot", "trow_sep"}
		for _, class := range omitClasses {
			if sel.HasClass(class) {
				return
			}
		}

		// Acá se podria subdividir la ejecución en hilos para procesar cada fila en paralelo.
		topic := s.parseThreadRow(sel)
		if topic.ID > 0 {
			exists := false
			for _, t := range result.Topics {
				if t.ID == topic.ID {
					exists = true
					break
				}
			}
			if !exists {
				result.Topics = append(result.Topics, topic)
			}
		}
	})

	result.Total = len(result.Topics)

	return result, nil
}

func (s *Scraper) parseThread(sel *goquery.Selection) models.Topic {
	topic := models.Topic{}

	linkSel := sel.Find("a[href*='tema-']")
	href, exists := linkSel.Attr("href")
	if !exists {
		return topic
	}

	topic.URL = href
	topic.Title = strings.TrimSpace(linkSel.Text())

	parts := strings.Split(href, "tid=")
	if len(parts) > 1 {
		topic.ID, _ = strconv.Atoi(strings.Split(parts[1], "&")[0])
	}

	authorSel := sel.Find(".author a, .author span a")
	topic.Author = strings.TrimSpace(authorSel.Text())
	topic.AuthorURL, _ = authorSel.Attr("href")

	viewsSel := sel.Find(".views")
	viewsStr := strings.ReplaceAll(viewsSel.Text(), ",", "")
	topic.Views, _ = strconv.Atoi(viewsStr)

	repliesSel := sel.Find(".replies")
	repliesStr := strings.ReplaceAll(repliesSel.Text(), ",", "")
	topic.Replies, _ = strconv.Atoi(repliesStr)

	return topic
}

func (s *Scraper) parseThreadRow(sel *goquery.Selection) models.Topic {
	topic := models.Topic{}

	linkSel := sel.Find("a[href*='tema-'][id]")
	href, exists := linkSel.Attr("href")
	if !exists {
		return topic
	}

	topic.URL, _ = url.QueryUnescape(href)
	topic.Title = strings.TrimSpace(linkSel.Text())

	idStr, _ := linkSel.Attr("id")
	partIdStr := strings.Split(idStr, "tid_")
	if len(partIdStr) > 1 {
		idStr = partIdStr[1]
	}
	if idStr != "" {
		topic.ID, _ = strconv.Atoi(idStr)
	}

	authorSel := sel.Find("td:nth-child(3) .author a")
	topic.Author = strings.TrimSpace(authorSel.First().Text())
	href, _ = authorSel.First().Attr("href")
	topic.AuthorURL, _ = url.QueryUnescape(href)

	// Materias y Aportes
	sel.Find("td:nth-child(4), td:nth-child(5)").Each(func(i int, td *goquery.Selection) {
		draftText, ok := td.Find("a").Attr("title")
		if !ok {
			return
		}

		listOfStr := strings.Split(strings.TrimSpace(draftText), "\n")

		switch i {
		case 0: // Aportes
			topic.Aportes = listOfStr
		case 1: // Materias
			topic.Materias = listOfStr
		}
	})

	// Respuestas, Vistas y Agradecimientos
	sel.Find("td:nth-child(n+6):nth-child(-n+8)").Each(func(i int, td *goquery.Selection) {
		textNum := strings.ReplaceAll(strings.TrimSpace(td.Text()), ".", "")
		switch i {
		case 0:
			topic.Replies, _ = strconv.Atoi(textNum)
		case 1:
			topic.Views, _ = strconv.Atoi(textNum)
		case 2:
			topic.ThankedCount, _ = strconv.Atoi(textNum)
		}
	})

	lastPost := &models.ForumLastPost{}
	lastPostSel := sel.Find(".lastposter")

	dateText := strings.TrimSpace(DateTextRegex.FindString(lastPostSel.Text()))
	if lastPostDate, err := time.Parse(TimeLayout, dateText); err == nil {
		lastPost.Date = &lastPostDate
	}

	lastPostSel.Find("a").Each(func(i int, a *goquery.Selection) {
		href, _ := a.Attr("href")
		href, _ = url.QueryUnescape(href)

		switch i {
		case 0:
			lastPost.URL = BaseURL + href
		case 1:
			lastPost.Author = strings.TrimSpace(a.Text())
			lastPost.AuthorURL = href
		}
	})

	topic.LastPost = lastPost
	return topic
}

func (s *Scraper) parseThreadFromLink(sel *goquery.Selection) models.Topic {
	topic := models.Topic{}

	href, exists := sel.Attr("href")
	if !exists {
		return topic
	}

	topic.URL = href
	topic.Title = strings.TrimSpace(sel.Text())

	idAttr, _ := sel.Attr("id")
	if idAttr != "" {
		parts := strings.Split(idAttr, "tid_")
		if len(parts) > 1 {
			topic.ID, _ = strconv.Atoi(parts[1])
		}
	}

	return topic
}

func (s *Scraper) GetTopic(slug string) (*models.TopicDetail, error) {
	fmt.Println("Fetching from : ", BaseURL+slug)
	doc, err := s.fetchDoc(BaseURL + slug)
	if err != nil {
		return nil, err
	}

	topicDetail := &models.TopicDetail{
		Posts: []models.Post{},
		Topic: models.Topic{
			URL: BaseURL + slug,
		},
	}

	// ID
	idInHeadElem, _ := doc.Find("head meta[name*='twitter:app:url:iphone']").Attr("content")
	urlWithID, err := url.Parse(idInHeadElem)
	if err == nil {
		idStr := urlWithID.Query().Get("tid")
		topicDetail.ID, _ = strconv.Atoi(idStr)
	}

	// Materias y Aportes
	materiasTags := doc.Find("span.xttag")
	materiasTags.Each(func(i int, span *goquery.Selection) {
		draftText, ok := span.Attr("title")
		if !ok {
			return
		}
		listOfStr := strings.Split(strings.TrimSpace(draftText), "\n")

		switch i {
		case 0: // Aportes
			topicDetail.Aportes = listOfStr
		case 1: // Materias
			topicDetail.Materias = listOfStr
		}
	})

	// Titulo del tema
	title := materiasTags.Parent().Text()
	endOfTitle := strings.Index(title, "\n")
	topicDetail.Title = strings.TrimSpace(title[:endOfTitle])

	// Autor del tema
	user := doc.Find("span.usuario_link[data-id] a").First()
	topicDetail.Author = strings.TrimSpace(user.Text())
	href, _ := user.Attr("href")
	topicDetail.AuthorURL, _ = url.QueryUnescape(href)

	// Posts
	postList := doc.Find("table.tborder[id*='post_']")
	postList.Each(func(i int, sel *goquery.Selection) {
		post := s.parsePost(sel)
		if post.ID > 0 || post.Author != "" {
			topicDetail.Posts = append(topicDetail.Posts, post)
		}
	})

	topicDetail.Replies = len(topicDetail.Posts) - 1

	return topicDetail, nil
}

func (s *Scraper) parsePost(sel *goquery.Selection) models.Post {
	post := models.Post{}

	// ID del post
	id := sel.AttrOr("id", "")
	if id != "" {
		parts := strings.Split(id, "post_")
		if len(parts) > 1 {
			post.ID, _ = strconv.Atoi(parts[1])
		}
	}

	// Autor del post (comentario o respuesta del primer topico)
	authorSel := sel.Find("span.usuario_link[data-id] a")
	post.Author = strings.TrimSpace(authorSel.First().Text())
	autorHref, _ := authorSel.First().Attr("href")
	post.AuthorURL, _ = url.QueryUnescape(autorHref)

	// Contenido del post (ReplyTo y el mensaje del comentario)
	contentSel := sel.Find("div.post_message")
	replyPost := contentSel.Find("blockquote cite a.quick_jump").First()
	if replyPost.Length() > 0 {
		replyHref, _ := replyPost.Attr("href")
		replyUrl, _ := url.Parse(replyHref)
		replyId := replyUrl.Query().Get("pid")
		fmt.Println("Reply to post ID: ", replyId)
		post.ReplyPostID, _ = strconv.Atoi(replyId)
	}

	contentSel.Find("blockquote").Remove() // Elimino el bloque de cita para quedarnos solo con el mensaje del post
	post.Content = strings.TrimSpace(contentSel.Text())

	contentSel.Find("a").Each(func(i int, a *goquery.Selection) {
		href, _ := a.Attr("href")
		href, _ = url.QueryUnescape(href) 
		if strings.Contains(href, "attachment.php?"){ // pequeña excepción donde los archivos adjuntos están dentro del mensaje.
			href = BaseURL + href
			post.Attachments = append(post.Attachments, href)
			return
		}
		post.Links = append(post.Links, href)
	})

	// Documentos e imagenes adjuntas al mensaje principal
	sel.Find("fieldset a[name='download'], fieldset img.attachment").Each(func(i int, sel *goquery.Selection) {
		if sel.Is("a") {
			href, _ := sel.Attr("href")
			href, _ = url.QueryUnescape(href)
			if strings.Contains(href, "attachment.php?aid=") {
				post.Attachments = append(post.Attachments, BaseURL + href)
			}
		}
		if sel.Is("img") {
			src, _ := sel.Attr("src")
			src, _ = url.QueryUnescape(src)
			if strings.Contains(src, "attachment.php?aid=") {
				post.Attachments = append(post.Attachments, BaseURL + src)
			}
		}
	})

	// Fecha del post
	timeSel := sel.Find("tbody tr:nth-child(2) td").First()
	timeStr := strings.TrimSpace(timeSel.Text())
	post.Date, _ = time.Parse(TimeLayout, timeStr)

	return post
}

func (s *Scraper) parseDate(dateStr string) time.Time {
	now := time.Now()
	dateStr = strings.ToLower(dateStr)

	if strings.Contains(dateStr, "hoy") {
		return now
	}
	if strings.Contains(dateStr, "ayer") {
		return now.AddDate(0, 0, -1)
	}

	formats := []string{
		"02/01/2006, 15:04",
		"02/01/2006 15:04",
		"02-01-2006 15:04",
		"02/01/2006",
		"02-01-2006",
	}

	for _, format := range formats {
		if t, err := time.Parse(format, dateStr); err == nil {
			return t
		}
	}

	return now
}

func (s *Scraper) Search(query string, forumID int, action string) (*models.SearchResult, error) {
	searchURL := BaseURL + "/search.php?"

	if action != "" {
		searchURL += "action=" + action
	} else {
		searchURL += "action=do_search"
	}

	if query != "" {
		searchURL += "&keywords=" + url.QueryEscape(query)
	}

	if forumID > 0 {
		searchURL += "&fids=" + strconv.Itoa(forumID)
	}

	doc, err := s.fetchDoc(searchURL)
	if err != nil {
		return nil, err
	}

	result := &models.SearchResult{
		Topics: []models.Topic{},
	}

	doc.Find(".thread, tr.thread").Each(func(_ int, sel *goquery.Selection) {
		topic := s.parseThread(sel)
		if topic.ID > 0 {
			result.Topics = append(result.Topics, topic)
		}
	})

	result.Total = len(result.Topics)

	return result, nil
}

func (s *Scraper) GetUser(username string) (*models.User, error) {
	doc, err := s.fetchDoc(BaseURL + "/usuario-" + username)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Username: username,
	}

	userField := doc.Find(".user_field")
	userField.Each(func(_ int, sel *goquery.Selection) {
		label := strings.TrimSpace(sel.Find(".user_field_label").Text())
		value := strings.TrimSpace(sel.Find(".user_field_value").Text())

		switch {
		case strings.Contains(label, "Unido"):
			user.JoinDate = value
		case strings.Contains(label, "Mensajes"):
			user.Posts, _ = strconv.Atoi(strings.ReplaceAll(value, ",", ""))
		case strings.Contains(label, "Ubicación"):
			user.Location = value
		case strings.Contains(label, "Sitio"):
			user.Website = value
		}
	})

	avatarImg := doc.Find(".user_avatar img")
	user.Avatar, _ = avatarImg.Attr("src")

	sigSel := doc.Find(".signature")
	user.Signature = strings.TrimSpace(sigSel.Text())

	return user, nil
}

func (s *Scraper) fetchDoc(url string) (*goquery.Document, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
	req.Header.Set("Accept-Language", "es-AR,es;q=0.9,en;q=0.8")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		return nil, fmt.Errorf("page not found: %s", url)
	}

	if resp.StatusCode == 403 {
		return nil, fmt.Errorf("access forbidden: %s", url)
	}

	return goquery.NewDocumentFromReader(resp.Body)
}
