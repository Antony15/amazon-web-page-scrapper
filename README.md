# sellerapp-go-test

REST API in Go to scrap the given amazon web page

There are two endpoints used here,
1. scrapurl :- This endpoint will scrap the following data from the page
  - Product Name/Title
  - Product image url
  - Product description
  - Product price
  - Product total number of reviews
Then will send an post to below endpoint with payload as the scrapped data. On success will return message "Requested URl Successfully scrapped & saved in database" with status code as 200.
 Sample Json Request :
 {
	"url": "https://www.amazon.com/PlayStation-4-Pro-1TB-Console/dp/B01LOP8EZC/"
 }
 
 2. writedocument :- This endpoint will upsert the post payload in mongodb. On success will return message "Request Wrote to document" with status code as 200.
 
