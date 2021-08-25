# Task 1 : Online Store

## Problem Description

We are members of the engineering team of an online store. When we look at ratings for our online store application, we received the following facts:

1. Customers were able to put items in their cart, check out, and then pay. After several days, many of our customers received calls from our Customer Service department stating that their orders have been canceled due to stock unavailability.
2. These bad reviews generally come within a week after our 12.12 event, in which we held a large flash sale and set up other major discounts to promote our store.

After checking in with our Customer Service and Order Processing departments, we received the following additional facts:
1. Our inventory quantities are often misreported, and some items even go as far as having a negative inventory quantity.
2. The misreported items are those that performed very well on our 12.12 event.
3. Because of these misreported inventory quantities, the Order Processing department was unable to fulfill a lot of orders, and thus requested help from our Customer Service department to call our customers and notify them that we have had to cancel their orders.

## Analysis

- The misreported inventory quantity might be caused by data races during concurrent modification on the data, ex: get product, save product, create order, cancel order
- The apparent indication can be seen from the resulting quality of inventory that can be in irrational state, such as some inventory having negative quantity.

## Solution

- This inconsistency can be prevented by using concurrency control during data modification.
- We will use database transaction and repeatable read isolation level to manage concurrent database operation while still maintain reasonable performance. 
- Product Entity will be the aggregate root for consistency boundary of Product Entity and Order Entity

